package camtest

import (
    "go.pirat.app/pi/camunda"
    "go.pirat.app/pi/camunda/worker"
    "log"
    "reflect"
    "encoding/json"
)

type ContextStub struct {
    *worker.ContextImpl
    failure bool
    error   bool

    vars camunda.Variables
}

func CreateTestContext(vars camunda.Variables) *ContextStub {
    task := &camunda.ResLockedExternalTask{
        Variables:   vars,
        BusinessKey: "test-business-key",
    }

    return &ContextStub{
        ContextImpl: worker.NewContext(nil, task, "test-context-worker"),
        vars:        vars,
        failure:     false,
        error:       false,
    }
}

func (w *ContextStub) Complete(tc *worker.TaskComplete) error {
    if tc.Variables == nil {
        return nil
    }

    for key, val := range tc.Variables {
        if val.Type == "Object" && reflect.TypeOf(val.Value).Kind() == reflect.String {
            m := make(map[string]interface{})
            _ = json.Unmarshal([]byte(val.Value.(string)), &m)
            val.Value = m
        }

        w.vars[key] = val
    }

    return nil
}

func (w *ContextStub) HandleFailure(query worker.TaskFailureRequest) error {
    log.Printf("Failure handler called: %s Details: %s \n", query.ErrorMessage, query.ErrorDetails)
    w.failure = true
    return nil
}

func (w ContextStub) HandleBPMNError(code int, message string) error {
    w.error = true

    log.Printf("BPMN error received: %d message: %s\n", code, message)

    return nil
}

func (w ContextStub) ExtendLock(_ string, _ int) error {
    return nil
}

func (w ContextStub) StartLockExtender() {
    return
}

func (w ContextStub) StopExtender() {
    return
}

func (w *ContextStub) Variables() camunda.Variables {
    return w.vars
}
