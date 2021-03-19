package camunda

import "fmt"

// ProcessDefinitionResponse a JSON object corresponding to the ProcessDefinition interface in the engine
type ProcessDefinitionResponse struct {
	// The id of the process definition
	ID string `json:"id"`
	// The key of the process definition, i.e., the id of the BPMN 2.0 XML process definition
	Key string `json:"key"`
	// The category of the process definition
	Category string `json:"category"`
	// The description of the process definition
	Description string `json:"description"`
	// The name of the process definition
	Name string `json:"name"`
	// The version of the process definition that the engine assigned to it
	Version int `json:"Version"`
	// The file name of the process definition
	Resource string `json:"resource"`
	// The deployment id of the process definition
	DeploymentID string `json:"deploymentId"`
	// The file name of the process definition diagram, if it exists
	Diagram string `json:"diagram"`
	// A flag indicating whether the definition is suspended or not
	Suspended bool `json:"suspended"`
	// The tenant id of the process definition
	TenantID string `json:"tenantId"`
	// The version tag of the process definition
	VersionTag string `json:"versionTag"`
	// History time to live value of the process definition. Is used within History cleanup
	HistoryTimeToLive int `json:"historyTimeToLive"`
	// A flag indicating whether the process definition is startable in Tasklist or not
	StartableInTaskList bool `json:"startableInTasklist"`
}

// ResActivityInstanceStatistics a JSON array containing statistics results per activity
type ResActivityInstanceStatistics struct {
	// The id of the activity the results are aggregated for
	Id string `json:"id"`
	// The total number of running instances of this activity
	Instances int `json:"instances"`
	// Number	The total number of failed jobs for the running instances.
	// Note: Will be 0 (not null), if failed jobs were excluded
	FailedJobs int `json:"failedJobs"`
	// Each item in the resulting array is an object which contains the following properties
	Incidents []ResActivityInstanceStatisticsIncident `json:"incidents"`
}

// ResInstanceStatistics a JSON array containing statistics results per process definition
type ResInstanceStatistics struct {
	// The id of the activity the results are aggregated for
	Id string `json:"id"`
	// The total number of running instances of this activity
	Instances int `json:"instances"`
	// Number	The total number of failed jobs for the running instances.
	// Note: Will be 0 (not null), if failed jobs were excluded
	FailedJobs int `json:"failedJobs"`
	// The process definition with the properties as described in the Get single definition method
	Definition ProcessDefinitionResponse `json:"definition"`
	// Each item in the resulting array is an object which contains the following properties
	Incidents []ResActivityInstanceStatisticsIncident `json:"incidents"`
}

// ResActivityInstanceStatisticsIncident a statistics incident
type ResActivityInstanceStatisticsIncident struct {
	// The type of the incident the number of incidents is aggregated for.
	// See the User Guide for a list of incident types
	IncidentType string `json:"incidentType"`
	// The total number of incidents for the corresponding incident type
	IncidentCount int `json:"incidentCount"`
}

// QueryProcessDefinitionBy path builder
type QueryProcessDefinitionBy struct {
	Id       string
	Key      string
	TenantId string
}

// ResGetStartFormKey a response from GetStartFormKey method
type ResGetStartFormKey struct {
	// The form key for the process definition
	Key string `json:"key"`
	// The context path of the process application
	ContextPath string `json:"contextPath"`
}

// ResBPMNProcessDefinition a JSON object containing the id of the definition and the BPMN 2.0 XML
type ResBPMNProcessDefinition struct {
	// The id of the process definition
	Id string `json:"id"`
	// An escaped XML string containing the XML that this definition was deployed with.
	// Carriage returns, line feeds and quotation marks are escaped
	Bpmn20Xml string `json:"bpmn20Xml"`
}

// Path a build path part
func (q *QueryProcessDefinitionBy) Path() string {
	if q.Key != "" && q.TenantId != "" {
		return fmt.Sprintf("key/%s/tenant-id/%s", q.Key, q.TenantId)
	} else if q.Key != "" {
		return fmt.Sprintf("key/%s", q.Key)
	}

	return q.Id
}

// ResStartedProcessDefinition ProcessDefinition for started
type ResStartedProcessDefinition struct {
	// The id of the process definition
	Id string `json:"id"`
	// The id of the process definition
	DefinitionId string `json:"definitionId"`
	// The business key of the process instance
	BusinessKey string `json:"businessKey"`
	// The case instance id of the process instance
	CaseInstanceId string `json:"caseInstanceId"`
	// The tenant id of the process instance
	TenantId string `json:"tenantId"`
	// A flag indicating whether the instance is still running or not
	Ended bool `json:"ended"`
	// A flag indicating whether the instance is suspended or not
	Suspended bool `json:"suspended"`
	// A JSON array containing links to interact with the instance
	Links []ResLink `json:"links"`
	// A JSON object containing a property for each of the latest variables
	Variables map[string]Variable `json:"variables"`
}

// ReqSubmitStartForm request a SubmitStartForm
type ReqSubmitStartForm struct {
	// A JSON object containing the variables the process is to be initialized with.
	// Each key corresponds to a variable name and each value to a variable value
	Variables map[string]Variable `json:"variables"`
	// A JSON object containing the business key the process is to be initialized with.
	// The business key uniquely identifies the process instance in the context of the given process definition
	BusinessKey string `json:"businessKey"`
}

// ReqSubmitStartForm response rrom SubmitStartForm method
type ResSubmitStartForm struct {
	Links        []ResLink `json:"links"`
	Id           string    `json:"id"`
	DefinitionId string    `json:"definitionId"`
	BusinessKey  string    `json:"businessKey"`
	Ended        bool      `json:"ended"`
	Suspended    bool      `json:"suspended"`
}

// ReqActivateOrSuspendById response ActivateOrSuspendById
type ReqActivateOrSuspendById struct {
	// A Boolean value which indicates whether to activate or suspend a given process definition. When the value
	// is set to true, the given process definition will be suspended and when the value is set to false,
	// the given process definition will be activated
	Suspended *bool `json:"suspended,omitempty"`
	// A Boolean value which indicates whether to activate or suspend also all process instances of the given process
	// definition. When the value is set to true, all process instances of the provided process definition will be
	// activated or suspended and when the value is set to false, the suspension state of all process instances of
	// the provided process definition will not be updated
	IncludeProcessInstances *bool `json:"includeProcessInstances,omitempty"`
	// The date on which the given process definition will be activated or suspended. If null, the suspension state
	// of the given process definition is updated immediately. The date must have the format yyyy-MM-dd'T'HH:mm:ss,
	// e.g., 2013-01-23T14:42:45
	ExecutionDate *Time `json:"executionDate,omitempty"`
}

// ReqActivateOrSuspendByKey response ActivateOrSuspendByKey
type ReqActivateOrSuspendByKey struct {
	// The key of the process definitions to activate or suspend
	ProcessDefinitionKey string `json:"processDefinitionKey"`
	// A Boolean value which indicates whether to activate or suspend a given process definition. When the value
	// is set to true, the given process definition will be suspended and when the value is set to false,
	// the given process definition will be activated
	Suspended *bool `json:"suspended,omitempty"`
	// A Boolean value which indicates whether to activate or suspend also all process instances of the given process
	// definition. When the value is set to true, all process instances of the provided process definition will be
	// activated or suspended and when the value is set to false, the suspension state of all process instances of
	// the provided process definition will not be updated
	IncludeProcessInstances *bool `json:"includeProcessInstances,omitempty"`
	// The date on which the given process definition will be activated or suspended. If null, the suspension state
	// of the given process definition is updated immediately. The date must have the format yyyy-MM-dd'T'HH:mm:ss,
	// e.g., 2013-01-23T14:42:45
	ExecutionDate *Time `json:"executionDate,omitempty"`
}
