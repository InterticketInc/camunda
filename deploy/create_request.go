package deploy

import "io"

// CreateRequest a request to deployment create
type CreateRequest struct {
	DeploymentName           string
	EnableDuplicateFiltering bool
	DeployChangedOnly        bool
	DeploymentSource         string
	TenantID                 string
	Resources                map[string]io.Reader
}
