package camunda

// StartInstanceRequest a JSON object with the following properties: (at least an empty JSON object {}
// or an empty request body)
type StartInstanceRequest struct {
	// A JSON object containing the variables the process is to be initialized with
	Variables map[string]*Variable `json:"variables,omitempty"`
	// The business key the process instance is to be initialized with.
	// The business key uniquely identifies the process instance in the context of the given process definition
	BusinessKey string `json:"businessKey,omitempty"`
	// The case instance id the process instance is to be initialized with
	CaseInstanceId string `json:"caseInstanceId,omitempty"`
	// Optional. A JSON array of instructions that specify which activities to start the process instance at.
	// If this property is omitted, the process instance starts at its default blank start event
	StartInstructions []*StartInstructionsRequest `json:"startInstructions,omitempty"`
	// Skip execution listener invocation for activities that are started or ended as part of this request
	// Note: This option is currently only respected when start instructions are submitted via
	// the startInstructions property
	SkipCustomListeners bool `json:"skipCustomListeners,omitempty"`
	// Skip execution of input/output variable mappings for activities that are started or ended as part of this request
	// Note: This option is currently only respected when start instructions are submitted via
	// the startInstructions property
	SkipIoMappings bool `json:"skipIoMappings,omitempty"`
	// Indicates if the variables, which was used by the process instance during execution, should be returned. Default value: false
	WithVariablesInReturn bool `json:"withVariablesInReturn,omitempty"`
}

// RestartInstanceRequest a request to restart instance
type RestartInstanceRequest struct {
	// A list of process instance ids to restart
	ProcessInstanceIds *string `json:"processInstanceIds,omitempty"`
	// A historic process instance query like the request body described by POST /history/process-instance
	HistoricProcessInstanceQuery *string `json:"historicProcessInstanceQuery,omitempty"`
	// Optional. A JSON array of instructions that specify which activities to start the process instance at.
	// If this property is omitted, the process instance starts at its default blank start event
	StartInstructions *[]StartInstructionsRequest `json:"startInstructions,omitempty"`
	// Skip execution listener invocation for activities that are started or ended as part of this request
	// Note: This option is currently only respected when start instructions are submitted via
	// the startInstructions property
	SkipCustomListeners *bool `json:"skipCustomListeners,omitempty"`
	// Skip execution of input/output variable mappings for activities that are started or ended as part of this request
	// Note: This option is currently only respected when start instructions are submitted via
	// the startInstructions property
	SkipIoMappings *bool `json:"skipIoMappings,omitempty"`
	// Set the initial set of variables during restart. By default, the last set of variables is used
	InitialVariables *bool `json:"initialVariables,omitempty"`
	// Do not take over the business key of the historic process instance.
	WithoutBusinessKey *bool `json:"withoutBusinessKey,omitempty"`
}

// StartInstructionsRequest a JSON array of instructions that specify which activities to start the process instance at
type StartInstructionsRequest struct {
	// Mandatory. One of the following values: startBeforeActivity, startAfterActivity, startTransition.
	// A startBeforeActivity instruction requests to start execution before entering a given activity.
	// A startAfterActivity instruction requests to start at the single outgoing sequence flow of a given activity.
	// A startTransition instruction requests to execute a specific sequence flow
	Type string `json:"type"`
	// Can be used with instructions of types startBeforeActivity and startAfterActivity.
	// Specifies the activity the instruction targets
	ActivityId *string `json:"activityId,omitempty"`
	// Can be used with instructions of types startTransition. Specifies the sequence flow to start
	TransitionId *string `json:"transitionId,omitempty"`
	// Can be used with instructions of type startBeforeActivity, startAfterActivity, and startTransition.
	// A JSON object containing variable key-value pairs
	Variables *map[string]VariableSet `json:"variables,omitempty"`
}
