package deploy

import "github.com/interticketinc/camunda"

// CreateResponse a JSON object corresponding to the DeploymentWithDefinitions interface in the engine
type CreateResponse struct {
	// The id of the deployment
	ID string `json:"id"`
	// The name of the deployment
	Name string `json:"name"`
	// The source of the deployment
	Source string `json:"source"`
	// The tenant id of the deployment
	TenantID string `json:"tenantId"`
	// The time when the deployment was created
	DeploymentTime camunda.Time `json:"deploymentTime"`
	// Link to the newly created deployment with method, href and rel
	Links []camunda.ResLink `json:"links"`
	// A JSON Object containing a property for each of the process definitions,
	// which are successfully deployed with that deployment
	DeployedProcessDefinitions map[string]camunda.ProcessDefinitionResponse `json:"deployedProcessDefinitions"`
	// A JSON Object containing a property for each of the case definitions,
	// which are successfully deployed with that deployment
	DeployedCaseDefinitions map[string]camunda.CaseDefinition `json:"deployedCaseDefinitions"`
	// A JSON Object containing a property for each of the decision definitions,
	// which are successfully deployed with that deployment
	DeployedDecisionDefinitions map[string]camunda.DecisionDefinition `json:"deployedDecisionDefinitions"`
	// A JSON Object containing a property for each of the decision requirements definitions,
	// which are successfully deployed with that deployment
	DeployedDecisionRequirementsDefinitions map[string]camunda.ResDecisionRequirementsDefinition `json:"deployedDecisionRequirementsDefinitions"`
}
