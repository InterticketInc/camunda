package camunda

// ProcessInstanceQuery query struct for process instance
type ProcessInstanceQuery struct {
    //ProcessInstanceIDs	Filter by a comma-separated list of process instance ids.
    ProcessInstanceIDs []string `url:"processInstanceIds,omitempty"`
    // BusinessKey	Filter by process instance business key.
    BusinessKey string `url:"businessKey,omitempty"`
    // BusinessKeyLike	Filter by process instance business key that the parameter is a substring of.
    BusinessKeyLike string `url:"businessKeyLike,omitempty"`
    // CaseInstanceID	Filter by case instance id.
    CaseInstanceID string `url:"caseInstanceId,omitempty"`
    //processDefinitionId	Filter by the process definition the instances run on.
    //processDefinitionKey	Filter by the key of the process definition the instances run on.
    //processDefinitionKeyIn	Filter by a comma-separated list of process definition keys. A process instance must have one of the given process definition keys.
    //processDefinitionKeyNotIn	Exclude instances by a comma-separated list of process definition keys. A process instance must not have one of the given process definition keys.
    //deploymentId	Filter by the deployment the id belongs to.
    //superProcessInstance	Restrict query to all process instances that are sub process instances of the given process instance. Takes a process instance id.
    //subProcessInstance	Restrict query to all process instances that have the given process instance as a sub process instance. Takes a process instance id.
    //superCaseInstance	Restrict query to all process instances that are sub process instances of the given case instance. Takes a case instance id.
    //subCaseInstance	Restrict query to all process instances that have the given case instance as a sub case instance. Takes a case instance id.
    //active	Only include active process instances. Value may only be true, as false is the default behavior.
    //suspended	Only include suspended process instances. Value may only be true, as false is the default behavior.
    //withIncident	Filter by presence of incidents. Selects only process instances that have an incident.
    //incidentId	Filter by the incident id.
    //incidentType	Filter by the incident type. See the User Guide for a list of incident types.
    //incidentMessage	Filter by the incident message. Exact match.
    //incidentMessageLike	Filter by the incident message that the parameter is a substring of.
    //tenantIdIn	Filter by a comma-separated list of tenant ids. A process instance must have one of the given tenant ids.
    //withoutTenantId	Only include process instances which belong to no tenant. Value may only be true, as false is the default behavior.
    //activityIdIn	Filter by a comma-separated list of activity ids. A process instance must currently wait in a leaf activity with one of the given activity ids.
    //rootProcessInstances	Restrict the query to all process instances that are top level process instances.
    //leafProcessInstances	Restrict the query to all process instances that are leaf instances. (i.e. don't have any sub instances)
    //processDefinitionWithoutTenantId	Only include process instances which process definition has no tenant id.
    //variables	Only include process instances that have variables with certain values. Variable filtering expressions are comma-separated and are structured as follows:
    //A valid parameter value has the form key_operator_value. key is the variable name, operator is the comparison operator to be used and value the variable value.
    //Note: Values are always treated as String objects on server side.
    //
    //Valid operator values are: eq - equal to; neq - not equal to; gt - greater than; gteq - greater than or equal to; lt - lower than; lteq - lower than or equal to; like.
    //key and value may not contain underscore or comma characters.
    //variableNamesIgnoreCase	Match all variable names in this query case-insensitively. If set to true variableName and variablename are treated as equal.
    //variableValuesIgnoreCase	Match all variable values in this query case-insensitively. If set to true variableValue and variablevalue are treated as equal.
    //sortBy	Sort the results lexicographically by a given criterion. Valid values are instanceId, definitionKey, definitionId, tenantId and businessKey. Must be used in conjunction with the sortOrder parameter.
    //sortOrder	Sort the results in a given order. Values may be asc for ascending order or desc for descending order. Must be used in conjunction with the sortBy parameter.
    //firstResult	Pagination of results. Specifies the index of the first result to return.
    //maxResults	Pagination of results. Specifies the maximum number of results to return. Will return less results if there are no more results left.
}

// ProcessInstance A JSON array of process instance objects. Each process instance object has the following properties:
type ProcessInstance struct {
    //id	String	The id of the process instance.
    ID string `json:"id"`
    //definitionId	String	The id of the process definition that this process instance belongs to.
    DefinitionID string `json:"definitionId,omitempty"`
    //businessKey	String	The business key of the process instance.
    BusinessKey string `json:"businessKey,omitempty"`
    //caseInstanceId	String	The id of the case instance associated with the process instance.
    CaseInstanceID string `json:"caseInstanceId,omitempty"`
    //suspended	Boolean	A flag indicating whether the process instance is suspended or not.
    Suspended bool `json:"suspended,omitempty"`
    //tenantId	String	The tenant id of the process instance.
    TenantID string `json:"tenantId,omitempty"`
}
