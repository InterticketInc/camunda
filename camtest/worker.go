package camtest

import (
	"github.com/interticketinc/camunda"
	"github.com/interticketinc/camunda/worker"
)

type ContextStub struct {
	*worker.ContextImpl
}

func CreateTestContext() *ContextStub {
	vars := camunda.Variables{}

	task := &camunda.ResLockedExternalTask{
		Variables:   vars,
		BusinessKey: "test-business-key",
	}

	return &ContextStub{
		ContextImpl: worker.NewContext(nil, task, "test-context-worker"),
	}
}

func (w ContextStub) Complete(tc *worker.TaskComplete) error {
	panic("implement me")
}

func (w ContextStub) HandleFailure(query worker.TaskFailureRequest) error {
	panic("implement me")
}

func (w ContextStub) HandleBPMNError(code int, message string) error {
	panic("implement me")
}

func (w ContextStub) ExtendLock(id string, duration int) error {
	panic("implement me")
}

func (w ContextStub) StartLockExtender() {
	panic("implement me")
}

func (w ContextStub) StopExtender() {
	panic("implement me")
}
