package camtest

import (
    "go.pirat.app/pi/camunda"
    "go.pirat.app/pi/camunda/worker"
    "reflect"
    "encoding/json"
    "errors"
)

type BpmnError struct {
    Code    int
    Message string
}

type ContextStub struct {
    *worker.ContextImpl

    vars   camunda.Variables
    failed *worker.TaskFailureRequest
    error  *BpmnError
}

func CreateTestContext(vars camunda.Variables) *ContextStub {
    task := &camunda.ResLockedExternalTask{
        Variables:   vars,
        BusinessKey: "test-business-key",
        TaskBase: &camunda.TaskBase{
            ID: "test-task-id",
        },
    }

    return &ContextStub{
        ContextImpl: worker.NewContext(nil, task, "test-context-worker"),
        vars:        vars,
    }
}

func (w *ContextStub) Complete(tc *worker.TaskComplete) error {
    if tc.Variables == nil {
        return nil
    }

    for key, val := range tc.Variables {
        if val.Type == "Object" && reflect.TypeOf(val.Value).Kind() == reflect.String {
            m := make(map[string]interface{})
            err := json.Unmarshal([]byte(val.Value.(string)), &m)

            if err != nil {
                var te *json.UnmarshalTypeError

                if errors.As(err, &te) && te.Value == "array" {
                    var l []interface{}
                    _ = json.Unmarshal([]byte(val.Value.(string)), &l)

                    val.ValueInfo = &camunda.ValueInfo{
                        ObjectTypeName:          "java.util.ArrayList",
                        SerializationDataFormat: "application/json",
                    }
                    val.Value = l
                }
            } else {
                val.Value = m
            }
        }

        w.vars[key] = val
    }

    return nil
}

func (w *ContextStub) HandleFailure(query worker.TaskFailureRequest) error {
    w.error = nil
    w.failed = &query

    return nil
}

func (w *ContextStub) HandleBPMNError(code int, message string) error {
    w.failed = nil
    w.error = &BpmnError{
        Code:    code,
        Message: message,
    }

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

func (w ContextStub) Failed() *worker.TaskFailureRequest {
    return w.failed
}

func (w ContextStub) Error() *BpmnError {
    return w.error
}
