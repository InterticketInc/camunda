package camunda

// TaskFilter is a filter array for filtering tasks tasks.
// The values will be marshalled into a query string
// More information at: https://docs.camunda.org/manual/7.14/reference/rest/task/get-query/
type TaskFilter struct {
	// ExternalTaskId filter by an external task's id.
	ExternalTaskId string `url:"externalTaskId,omitempty"`
	// ExternalTaskIDIn Filter by the comma-separated list of external task ids.
	ExternalTaskIDIn string `url:"externalTaskIDIn,omitempty"`
	// TopicName Filter by an external task topic.
	TopicName string `url:"topicName,omitempty"`
	// WorkerID Filter by the id of the worker that the task was most recently Locked by.
	WorkerID string `url:"workerId,omitempty"`
	// Locked Only include external tasks that are currently Locked (i.e., they have a lock time and it has not expired). Value may only be true, as false matches any external task.
	Locked string `url:"locked,omitempty"`
	// NotLocked Only include external tasks that are currently not Locked (i.e., they have no lock or it has expired). Value may only be true, as false matches any external task.
	NotLocked string `url:"notLocked,omitempty"`
	// WithRetriesLeft Only include external tasks that have a positive (> 0) number of retries (or null). Value may only be true, as false matches any external task.
	WithRetriesLeft string `url:"withRetriesLeft,omitempty"`
	// NoRetriesLeft Only include external tasks that have 0 retries. Value may only be true, as false matches any external task.
	NoRetriesLeft string `url:"noRetriesLeft,omitempty"`
	// LockExpirationAfter Restrict to external tasks that have a lock that expires after a given date. By default*, the date must have the format yyyy-MM-dd'T'HH:mm:ss.SSSZ, e.g., 2013-01-23T14:42:45.000+0200.
	LockExpirationAfter string `url:"lockExpirationAfter,omitempty"`
	// LockExpirationBefore Restrict to external tasks that have a lock that expires before a given date. By default*, the date must have the format yyyy-MM-dd'T'HH:mm:ss.SSSZ, e.g., 2013-01-23T14:42:45.000+0200.
	LockExpirationBefore string `url:"lockExpirationBefore,omitempty"`
	// ActivityID Filter by the id of the activity that an external task is created for.
	ActivityID string `url:"activityId,omitempty"`
	// ActivityIDIn Filter by the comma-separated list of ids of the activities that an external task is created for.
	ActivityIDIn string `url:"activityIdIn,omitempty"`
	// ExecutionId Filter by the id of the execution that an external task belongs to.
	ExecutionId string `url:"executionId,omitempty"`
	// ProcessInstanceID Filter by the id of the process instance that an external task belongs to.
	ProcessInstanceID string `url:"processInstanceId,omitempty"`
	// ProcessDefinitionID Filter by the id of the process definition that an external task belongs to.
	ProcessDefinitionID string `url:"processDefinitionId,omitempty"`
	// TenantIDIn Filter by a comma-separated list of tenant ids. An external task must have one of the given tenant ids.
	TenantIDIn []string `url:"tenantIdIn,omitempty"`
	// Active Only include active tasks. Value may only be true, as false matches any external task.
	Active string `url:"active,omitempty"`
	// PriorityHigherThanOrEquals Only include jobs with a priority higher than or equal to the given value. Value must be a valid long value.
	PriorityHigherThanOrEquals string `url:"priorityHigherThanOrEquals,omitempty"`
	// PriorityLowerThanOrEquals Only include jobs with a priority lower than or equal to the given value. Value must be a valid long value.
	PriorityLowerThanOrEquals string `url:"priorityLowerThanOrEquals,omitempty"`
	// Suspended Only include suspended tasks. Value may only be true, as false matches any external task.
	Suspended string `url:"suspended,omitempty"`
	// SortBy Sort the results lexicographically by a given criterion. Valid values are id, lockExpirationTime, processInstanceId, processDefinitionId, processDefinitionKey, tenantId and taskPriority. Must be used in conjunction with the sortOrder parameter.
	SortBy string `url:"sortBy,omitempty"`
	// SortOrder Sort the results in a given order. Values may be asc for ascending order or desc for descending order. Must be used in conjunction with the sortBy parameter.
	SortOrder string `url:"sortOrder,omitempty"`
	// FirstResult Pagination of results. Specifies the index of the first result to return.
	FirstResult string `url:"firstResult,omitempty"`
	// MaxResults Pagination of results. Specifies the maximum number of results to return. Will return less results if there are no more results left.
	MaxResults string `url:"maxResults,omitempty"`
}
