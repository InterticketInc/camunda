package worker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/interticketinc/camunda"
	"strconv"
)

// extendDuration in seconds how much should the extender increase the locking
const extendDuration = 10

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
type Handler func(ctx Context) error

type Context interface {
	Complete(tc *TaskComplete) error
	HandleFailure(query TaskFailureRequest) error
	HandleBPMNError(code int, message string) error
	ExtendLock(id string, duration int) error
	Variables() camunda.Variables
	StartLockExtender()
	StopExtender()
	TaskID() string
	TopicName() string
	Retries() int
}

// NewContext creates a new worker task ContextImpl
func NewContext(client *camunda.Client, task *camunda.ResLockedExternalTask, workerID string) *ContextImpl {
	return &ContextImpl{
		Task:     task,
		client:   client,
		workerID: workerID,
	}
}

// ContextImpl external task ContextImpl
type ContextImpl struct {
	Task     *camunda.ResLockedExternalTask
	client   *camunda.Client
	workerID string

	// Extender stop channel
	done chan interface{}
}

func (c *ContextImpl) TaskID() string {
	return c.Task.ID
}

func (c *ContextImpl) Variables() camunda.Variables {
	return c.Task.Variables
}

func (c *ContextImpl) TopicName() string {
	return c.Task.TopicName
}

func (c *ContextImpl) Retries() int {
	return c.Task.Retries
}

// Complete a mark external task is complete
func (c *ContextImpl) Complete(tc *TaskComplete) error {
	tm := c.client.TaskManager()
	return tm.Complete(c.Task.ID, camunda.QueryComplete{
		WorkerID:       &c.Task.WorkerID,
		Variables:      tc.Variables,
		LocalVariables: tc.LocalVariables,
	})
}

// HandleFailure handle external task failure
func (c *ContextImpl) HandleFailure(query TaskFailureRequest) error {
	return c.client.TaskManager().TaskFailed(c.Task.ID, camunda.Failure{
		WorkerID:     c.Task.WorkerID,
		ErrorMessage: query.ErrorMessage,
		ErrorDetails: query.ErrorDetails,
		Retries:      query.Retries,
		RetryTimeout: query.RetryTimeout,
	})
}

// HandleBPMNError handle external task failure
func (c *ContextImpl) HandleBPMNError(code int, message string) error {
	return c.client.TaskManager().HandleBPMNError(c.Task.ID, camunda.QueryHandleBPMNError{
		WorkerID:     c.Task.WorkerID,
		ErrorMessage: message,
		ErrorCode:    strconv.Itoa(code),
	})
}

// ExtendLock extending lock on specific task ID
func (c *ContextImpl) ExtendLock(id string, duration int) error {
	tm := c.client.TaskManager()
	err := tm.ExtendLock(id, camunda.QueryExtendLock{
		NewDuration: duration,
		WorkerID:    c.workerID,
	})
	if err != nil {
		return fmt.Errorf("error while extending lock: %w", err)
	}

	return nil
}

// StartLockExtender automatically extends a lock on a remote task if the task
// exceeds the default locking time
// You should call StopExtender() to clean up resources
func (c *ContextImpl) StartLockExtender() {
	c.done = make(chan interface{}, 1)

	go func() {
		t := time.NewTicker(9 * time.Second)

		for {
			select {
			case <-t.C:
				log.Debug().
					Str("task", c.Task.ID).
					Msg("extending lock on task")

				err := c.ExtendLock(c.Task.ID, extendDuration*1000)
				if err != nil {
					log.Err(err).Msg("failed to extend lock")
					return
				}
			case <-c.done:
				return
			}
		}
	}()
}

// StopExtender stops the task lock extender
func (c *ContextImpl) StopExtender() {
	if c.done != nil {
		// Stopping the extender goroutine
		c.done <- struct{}{}
		close(c.done)
	}
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
		ctx := NewContext(p.client, task, p.options.WorkerID)
		p.handle(ctx, handler)
	}
}

func (p *Worker) handle(ctx Context, handler Handler) {
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
	}
}
