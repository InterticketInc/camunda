package worker

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"go.pirat.app/pi/camunda"
)

// Worker external task worker
type Worker struct {
	client  *camunda.Client
	options *Options
	logger  func(err error)
}

// Options options for Worker
type Options struct {
	// WorkerID for all request (default: `worker-{random_int}`)
	WorkerID string `json:"workerId"`
	// LockDuration lock duration for all external task
	LockDuration time.Duration
	// MaxTasks maximum tasks to receive for 1 request to camunda
	MaxTasks int
	// MaxParallelTaskPerHandler maximum running parallel task per handler
	MaxParallelTaskPerHandler int
	// UsePriority use priority
	UsePriority *bool
	// LongPollingTimeout long polling timeout
	LongPollingTimeout time.Duration
}

// New a create new instance Worker
func New(client *camunda.Client, options *Options, logger func(err error)) *Worker {
	if options.WorkerID == "" {
		rand.Seed(time.Now().UnixNano())
		options.WorkerID = fmt.Sprintf("worker-%d", rand.Int())
	}

	return &Worker{
		client:  client,
		options: options,
		logger:  logger,
	}
}

// Handler a handler for external task
type Handler func(ctx *Context) error

// Context external task context
type Context struct {
	Task   *camunda.ResLockedExternalTask
	client *camunda.Client
}

// Complete a mark external task is complete
func (c *Context) Complete(query CompleteRequest) error {
	tm := c.client.TaskManager()
	return tm.Complete(c.Task.ID, camunda.QueryComplete{
		WorkerID:       &c.Task.WorkerID,
		Variables:      query.Variables,
		LocalVariables: query.LocalVariables,
	})
}

// HandleBPMNError handle external task BPMN error
func (c *Context) HandleBPMNError(query QueryHandleBPMNError) error {
	return c.client.TaskManager().HandleBPMNError(c.Task.ID, camunda.QueryHandleBPMNError{
		WorkerID:     &c.Task.WorkerID,
		ErrorCode:    query.ErrorCode,
		ErrorMessage: query.ErrorMessage,
		Variables:    query.Variables,
	})
}

// HandleFailure handle external task failure
func (c *Context) HandleFailure(query QueryHandleFailure) error {
	return c.client.TaskManager().HandleFailure(c.Task.ID, camunda.QueryHandleFailure{
		WorkerID:     &c.Task.WorkerID,
		ErrorMessage: query.ErrorMessage,
		ErrorDetails: query.ErrorDetails,
		Retries:      query.Retries,
		RetryTimeout: query.RetryTimeout,
	})
}

// AddHandler a add handler for external task
func (p *Worker) AddHandler(topics *[]camunda.QueryFetchAndLockTopic, handler Handler) {
	if topics != nil && p.options.LockDuration != 0 {
		for i := range *topics {
			v := &(*topics)[i]

			if v.LockDuration <= 0 {
				v.LockDuration = int(p.options.LockDuration / time.Millisecond)
			}
		}
	}

	var asyncResponseTimeout *int
	msValue := int(p.options.LongPollingTimeout.Nanoseconds() / int64(time.Millisecond))
	asyncResponseTimeout = &msValue

	go p.startPuller(camunda.QueryFetchAndLock{
		WorkerID:             p.options.WorkerID,
		MaxTasks:             p.options.MaxTasks,
		UsePriority:          p.options.UsePriority,
		AsyncResponseTimeout: asyncResponseTimeout,
		Topics:               topics,
	}, handler)
}

func (p *Worker) startPuller(query camunda.QueryFetchAndLock, handler Handler) {
	tasksChan := make(chan *camunda.ResLockedExternalTask)

	maxParallelTaskPerHandler := p.options.MaxParallelTaskPerHandler
	if maxParallelTaskPerHandler < 1 {
		maxParallelTaskPerHandler = 1
	}

	// create worker pool
	for i := 0; i < maxParallelTaskPerHandler; i++ {
		go p.runWorker(handler, tasksChan)
	}

	retries := 0
	for {
		tasks, err := p.client.TaskManager().FetchAndLock(query)
		if err != nil {
			if retries < 60 {
				retries++
			}

			p.logger(fmt.Errorf("failed pull: %w, sleeping: %d seconds", err, retries))
			time.Sleep(time.Duration(retries) * time.Second)
			continue
		}
		retries = 0

		for _, task := range tasks {
			tasksChan <- task
		}
	}
}

func (p *Worker) runWorker(handler Handler, tasksChan chan *camunda.ResLockedExternalTask) {
	for task := range tasksChan {
		p.handle(&Context{
			Task:   task,
			client: p.client,
		}, handler)
	}
}

func (p *Worker) handle(ctx *Context, handler Handler) {
	defer func() {
		if r := recover(); r != nil {
			errMessage := fmt.Sprintf("fatal error in task: %s", r)
			errDetails := fmt.Sprintf("fatal error in task: %s\nStack trace: %s", r, string(debug.Stack()))
			err := ctx.HandleFailure(QueryHandleFailure{
				ErrorMessage: &errMessage,
				ErrorDetails: &errDetails,
			})
			if err != nil {
				p.logger(fmt.Errorf("error send handle failure: %w", err))
			}

			p.logger(errors.New(errDetails))
		}
	}()

	err := handler(ctx)
	if err != nil {
		errMessage := fmt.Sprintf("task error: %s", err)
		err = ctx.HandleFailure(QueryHandleFailure{
			ErrorMessage: &errMessage,
		})

		if err != nil {
			p.logger(fmt.Errorf("error send handle failure: %w", err))
		}

		p.logger(errors.New(errMessage))
	}
}
