package deploy

// ResourceResponse a JSON array containing all deployment resources of the given deployment
type ResourceResponse struct {
	// The id of the deployment resource
	ID string `json:"id"`
	// The name of the deployment resource
	Name string `json:"name"`
	// The id of the deployment
	DeploymentID string `json:"deploymentId"`
}
