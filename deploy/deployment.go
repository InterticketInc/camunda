package deploy

import "github.com/interticketinc/camunda"

// Deployment a JSON array of deployment objects
type Deployment struct {
	// The id of the deployment
	ID string `json:"id"`
	// The name of the deployment
	Name string `json:"name"`
	// The source of the deployment
	Source string `json:"source"`
	// The tenant id of the deployment
	TenantID string `json:"tenantId"`
	// The date and time of the deployment.
	DeploymentTime camunda.Time `json:"deploymentTime"`
}
