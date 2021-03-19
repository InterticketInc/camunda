package camunda

import "fmt"

// TaskManager a client for ExternalTask API
type TaskManager struct {
	client *Client
}

// Get retrieves an external task by id, corresponding to the TaskManager interface in the engine
func (e *TaskManager) Get(id string) (*ResExternalTask, error) {
	resp := &ResExternalTask{}
	res, err := e.client.Get(
		"/external-task/"+id,
		map[string]string{},
	)
	if err != nil {
		return nil, err
	}

	if err := e.client.Marshal(res, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetList queries for the external tasks that fulfill given parameters.
// Parameters may be static as well as dynamic runtime properties of executions
// Query parameters described in the documentation:
// https://docs.camunda.org/manual/latest/reference/rest/external-task/get-query/#query-parameters
func (e *TaskManager) GetList(filter *TaskFilter) ([]*ResExternalTask, error) {
	var resp []*ResExternalTask
	res, err := e.client.Get(
		"/external-task",
		filter,
	)
	if err != nil {
		return nil, err
	}

	if err := e.client.Marshal(res, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetListCount queries for the number of external tasks that fulfill given parameters.
// Takes the same parameters as the Get External Tasks method.
// Query parameters described in the documentation:
// https://docs.camunda.org/manual/latest/reference/rest/external-task/get-query-count/#query-parameters
func (e *TaskManager) GetListCount(query map[string]string) (int, error) {
	resCount := ResponseCount{}
	res, err := e.client.Get("/external-task/count", query)
	if err != nil {
		return 0, err
	}

	err = e.client.Marshal(res, &resCount)
	return resCount.Count, err
}

// GetListPost queries for external tasks that fulfill given parameters in the form of a JSON object.
// This method is slightly more powerful than the Get External Tasks method
// because it allows to specify a hierarchical result sorting.
func (e *TaskManager) GetListPost(query QueryGetListPost, firstResult, maxResults int) ([]*ResExternalTask, error) {
	resp := []*ResExternalTask{}
	res, err := e.client.Post(
		"/external-task",
		map[string]string{},
		&query,
	)
	if err != nil {
		return nil, err
	}

	if err := e.client.Marshal(res, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// GetListPostCount queries for the number of external tasks that fulfill given parameters.
// This method takes the same message body as the Get External Tasks (POST) method
func (e *TaskManager) GetListPostCount(query QueryGetListPost) (int, error) {
	resCount := ResponseCount{}
	res, err := e.client.Post(
		"/external-task/count",
		map[string]string{},
		query,
	)
	if err != nil {
		return 0, err
	}

	err = e.client.Marshal(res, resCount)
	return resCount.Count, err
}

// FetchAndLock fetches and locks a specific number of external tasks for execution by a worker.
// Query can be restricted to specific task topics and for each task topic an individual lock time can be provided
func (e *TaskManager) FetchAndLock(query QueryFetchAndLock) ([]*ResLockedExternalTask, error) {
	var resp []*ResLockedExternalTask
	res, err := e.client.Post(
		"/external-task/fetchAndLock",
		map[string]string{},
		&query,
	)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	if err := e.client.Marshal(res, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Complete a completes an external task by id and updates process variables
func (e *TaskManager) Complete(id string, query QueryComplete) error {
	_, err := e.client.Post("/external-task/"+id+"/complete", map[string]string{}, &query)
	return err
}

// HandleBPMNError reports a business error in the context of a running external task by id.
// The error code must be specified to identify the BPMN error handler
func (e *TaskManager) HandleBPMNError(id string, query QueryHandleBPMNError) error {
	_, err := e.client.Post("/external-task/"+id+"/bpmnError", map[string]string{}, &query)
	return err
}

// HandleFailure reports a failure to execute an external task by id.
// A number of retries and a timeout until the task can be retried can be specified.
// If retries are set to 0, an incident for this task is created
func (e *TaskManager) HandleFailure(id string, query QueryHandleFailure) error {
	_, err := e.client.Post("/external-task/"+id+"/failure", map[string]string{}, &query)
	return err
}

// Unlock a unlocks an external task by id. Clears the task’s lock expiration time and worker id
func (e *TaskManager) Unlock(id string) error {
	_, err := e.client.doPost("/external-task/"+id+"/unlock", map[string]string{})
	return err
}

// ExtendLock a extends the timeout of the lock by a given amount of time
func (e *TaskManager) ExtendLock(id string, query QueryExtendLock) error {
	_, err := e.client.Post("/external-task/"+id+"/extendLock", map[string]string{}, &query)
	return err
}

// SetPriority a sets the priority of an existing external task by id. The default value of a priority is 0
func (e *TaskManager) SetPriority(id string, priority int) error {
	_, err := e.client.doPut("/external-task/"+id+"/priority", map[string]string{})
	return err
}

// SetRetries a sets the number of retries left to execute an external task by id. If retries are set to 0,
// an incident is created
func (e *TaskManager) SetRetries(id string, retries int) error {
	return e.client.doPutJSON("/external-task/"+id+"/retries", map[string]string{}, map[string]int{
		"retries": retries,
	})
}

// SetRetriesAsync a set Retries For Multiple External Tasks Async (Batch).
// Sets the number of retries left to execute external tasks by id asynchronously.
// If retries are set to 0, an incident is created
func (e *TaskManager) SetRetriesAsync(id string, query QuerySetRetriesAsync) (*ResBatch, error) {
	resp := ResBatch{}
	res, err := e.client.Post(
		"/external-task/retries-async",
		map[string]string{},
		&query,
	)
	if err != nil {
		return nil, err
	}

	if err := e.client.Marshal(res, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetRetriesSync a set Retries For Multiple External Tasks Sync.
// Sets the number of retries left to execute external tasks by id synchronously.
// If retries are set to 0, an incident is created
func (e *TaskManager) SetRetriesSync(id string, query QuerySetRetriesSync) error {
	return e.client.doPutJSON("/external-task/retries", map[string]string{}, &query)
}
