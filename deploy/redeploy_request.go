package deploy

// RedeployRequest a request to redeploy
type RedeployRequest struct {
	// A list of deployment resource ids to re-deploy
	ResourceIds *string `json:"resourceIds,omitempty"`
	// A list of deployment resource names to re-deploy
	ResourceNames *string `json:"resourceNames,omitempty"`
	// Sets the source of the deployment
	Source *string `json:"source,omitempty"`
}
