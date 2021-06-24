package worker

import "go.pirat.app/pi/camunda"

// TaskComplete a query for Complete request
type TaskComplete struct {
	// A JSON object containing variable key-value pairs
	Variables camunda.Variables `json:"variables"`
	// A JSON object containing variable key-value pairs.
	// Local variables are set only in the scope of external task
	LocalVariables camunda.Variables `json:"localVariables"`
}

// TaskFailureRequest a query for TaskFailed request
type TaskFailureRequest struct {
	// An message indicating the reason of the failure
	ErrorMessage string `json:"errorMessage,omitempty"`
	// A detailed error description
	ErrorDetails string `json:"errorDetails,omitempty"`
	// A number of how often the task should be retried.
	// Must be >= 0. If this is 0, an incident is created and the task cannot be fetched anymore unless
	// the retries are increased again. The incident's message is set to the errorMessage parameter
	Retries int `json:"retries,omitempty"`
	// A timeout in milliseconds before the external task becomes available again for fetching. Must be >= 0
	RetryTimeout int `json:"retryTimeout,omitempty"`
}
