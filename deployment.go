package camunda

// ResDeployment a JSON array of deployment objects
type ResDeployment struct {
	// The id of the deployment
	ID string `json:"id"`
	// The name of the deployment
	Name string `json:"name"`
	// The source of the deployment
	Source string `json:"source"`
	// The tenant id of the deployment
	TenantID string `json:"tenantId"`
	// The date and time of the deployment.
	DeploymentTime Time `json:"deploymentTime"`
}

// ResDeploymentCreate a JSON object corresponding to the DeploymentWithDefinitions interface in the engine
type ResDeploymentCreate struct {
	// The id of the deployment
	ID string `json:"id"`
	// The name of the deployment
	Name string `json:"name"`
	// The source of the deployment
	Source string `json:"source"`
	// The tenant id of the deployment
	TenantID string `json:"tenant_id"`
	// The time when the deployment was created
	DeploymentTime Time `json:"deployment_time"`
	// Link to the newly created deployment with method, href and rel
	Links []ResLink `json:"links"`
	// A JSON Object containing a property for each of the process definitions,
	// which are successfully deployed with that deployment
	DeployedProcessDefinitions map[string]ResProcessDefinition `json:"deployedProcessDefinitions"`
	// A JSON Object containing a property for each of the case definitions,
	// which are successfully deployed with that deployment
	DeployedCaseDefinitions map[string]ResCaseDefinition `json:"deployedCaseDefinitions"`
	// A JSON Object containing a property for each of the decision definitions,
	// which are successfully deployed with that deployment
	DeployedDecisionDefinitions map[string]*DecisionDefinition `json:"deployedDecisionDefinitions"`
	// A JSON Object containing a property for each of the decision requirements definitions,
	// which are successfully deployed with that deployment
	DeployedDecisionRequirementsDefinitions map[string]ResDecisionRequirementsDefinition `json:"deployedDecisionRequirementsDefinitions"`
}

// ReqDeploymentCreate a request to deployment create
type ReqDeploymentCreate struct {
	DeploymentName           string
	EnableDuplicateFiltering *bool
	DeployChangedOnly        *bool
	DeploymentSource         *string
	TenantID                 *string
	Resources                map[string]interface{}
}

// ReqRedeploy a request to redeploy
type ReqRedeploy struct {
	// A list of deployment resource ids to re-deploy
	ResourceIds *string `json:"resourceIds,omitempty"`
	// A list of deployment resource names to re-deploy
	ResourceNames *string `json:"resourceNames,omitempty"`
	// Sets the source of the deployment
	Source *string `json:"source,omitempty"`
}

// ResDeploymentResource a JSON array containing all deployment resources of the given deployment
type ResDeploymentResource struct {
	// The id of the deployment resource
	ID string `json:"id"`
	// The name of the deployment resource
	Name string `json:"name"`
	// The id of the deployment
	DeploymentID string `json:"deploymentId"`
}
