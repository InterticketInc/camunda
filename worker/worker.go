package worker

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"runtime/debug"
	"time"

	"go.pirat.app/pi/camunda"
)

// Worker external task worker
type Worker struct {
	client  *camunda.Client
	options *Options
	log     zerolog.Logger
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
func New(client *camunda.Client, options *Options) *Worker {
	if options.WorkerID == "" {
		rand.Seed(time.Now().UnixNano())
		options.WorkerID = fmt.Sprintf("worker-%d", rand.Int())
	}

	return &Worker{
		client:  client,
		options: options,
		log: log.With().
			Caller().
			Str("worker", options.WorkerID).
			Logger(),
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
func (c *Context) HandleBPMNError(query camunda.QueryHandleBPMNError) error {
	return c.client.TaskManager().HandleBPMNError(c.Task.ID, camunda.QueryHandleBPMNError{
		WorkerID:     &c.Task.WorkerID,
		ErrorCode:    query.ErrorCode,
		ErrorMessage: query.ErrorMessage,
		Variables:    query.Variables,
	})
}

// HandleFailure handle external task failure
func (c *Context) HandleFailure(query TaskFailureRequest) error {
	return c.client.TaskManager().TaskFailed(c.Task.ID, camunda.Failure{
		WorkerID:     c.Task.WorkerID,
		ErrorMessage: query.ErrorMessage,
		ErrorDetails: query.ErrorDetails,
		Retries:      query.Retries,
		RetryTimeout: query.RetryTimeout,
	})
}

// AddHandler a add handler for external task
func (p *Worker) AddHandler(topics []*camunda.TopicLockConfig, handler Handler) {
	if topics != nil && p.options.LockDuration != 0 {
		for i := range topics {
			v := topics[i]

			if v.LockDuration <= 0 {
				v.LockDuration = int(p.options.LockDuration / time.Millisecond)
			}
		}
	}

	var asyncResponseTimeout *int
	msValue := int(p.options.LongPollingTimeout.Nanoseconds() / int64(time.Millisecond))
	asyncResponseTimeout = &msValue

	go p.startPuller(camunda.FetchAndLockRequest{
		WorkerID:             p.options.WorkerID,
		MaxTasks:             p.options.MaxTasks,
		UsePriority:          p.options.UsePriority,
		AsyncResponseTimeout: asyncResponseTimeout,
		Topics:               topics,
	}, handler)
}

func (p *Worker) startPuller(req camunda.FetchAndLockRequest, handler Handler) {
	tasksChan := make(chan *camunda.ResLockedExternalTask)

	maxParallelTaskPerHandler := p.options.MaxParallelTaskPerHandler
	if maxParallelTaskPerHandler < 1 {
		maxParallelTaskPerHandler = 1
	}

	// create worker pool
	for i := 0; i < maxParallelTaskPerHandler; i++ {
		go p.runWorker(handler, tasksChan)
	}

	delay := 0
	for {
		tasks, err := p.client.TaskManager().FetchAndLock(req)
		if err != nil {
			if delay < 60 {
				delay++
			}

			bb, _ := json.Marshal(req)

			p.log.Error().Err(err).
				RawJSON("req", bb).
				Msgf("failed to pull message! sleeping: %d seconds", delay)
			time.Sleep(time.Duration(delay) * time.Second)
			continue
		}
		delay = 0

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
			err := ctx.HandleFailure(TaskFailureRequest{
				ErrorMessage: errMessage,
				ErrorDetails: errDetails,
			})
			if err != nil {
				p.log.Error().
					Err(err).
					Msg("error send handle failure")
			}

			p.log.Error().Msg(errMessage)
		}
	}()

	err := handler(ctx)
	if err != nil {
		msg := fmt.Sprintf("task error: %s", err)
		err = ctx.HandleFailure(TaskFailureRequest{
			ErrorMessage: msg,
		})

		if err != nil {
			p.log.Error().
				Err(err).
				Msg("error send handle failure")
		}

		p.log.Error().Msg(msg)
	}
}
