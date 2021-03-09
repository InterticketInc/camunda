package camunda

type TaskBase struct {
	// The id of the activity that this external task belongs to
	ActivityId string `json:"activityId"`
	// The id of the activity instance that the external task belongs to
	ActivityInstanceId string `json:"activityInstanceId"`
	// The full error message submitted with the latest reported failure executing this task;
	// null if no failure was reported previously or if no error message was submitted
	ErrorMessage string `json:"errorMessage"`
	// The error message that was supplied when the last failure of this task was reported.
	ErrorDetails string `json:"errorDetails"`
	// The id of the execution that the external task belongs to.
	ExecutionId string `json:"executionId"`
	// The id of the external task.
	Id string `json:"id"`
	// The date that the task's most recent lock expires or has expired
	LockExpirationTime string `json:"lockExpirationTime"`
	// The id of the process definition the external task is defined in
	ProcessDefinitionId string `json:"processDefinitionId"`
	// The key of the process definition the external task is defined in
	ProcessDefinitionKey string `json:"processDefinitionKey"`
	// The id of the process instance the external task belongs to
	ProcessInstanceId string `json:"processInstanceId"`
	// The id of the tenant the external task belongs to
	TenantId string `json:"tenantId"`
	// The number of retries the task currently has left
	Retries int `json:"retries"`
	// A flag indicating whether the external task is suspended or not
	Suspended bool `json:"suspended"`
	// The id of the worker that possesses or possessed the most recent lock
	WorkerID string `json:"workerId"`
	// The priority of the external task
	Priority int `json:"priority"`
	// The topic name of the external task
	TopicName string `json:"topicName"`
}