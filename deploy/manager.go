package deploy

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"

	"go.pirat.app/pi/camunda"
)

// Manager deployment manager instance. You can instantiate the instance from the camunda.Client instance
type Manager struct {
	client *camunda.Client
}

// NewManager initializes a new deployment manager
func NewManager(client *camunda.Client) *Manager {
	return &Manager{
		client: client,
	}
}

// GetList a queries for deployments that fulfill given parameters. Parameters may be the properties of deployments,
// such as the id or name or a range of the deployment time. The size of the result set can be retrieved by using
// the Get Deployment count method.
// Query parameters described in the documentation:
// https://docs.camunda.org/manual/latest/reference/rest/deployment/get-query/#query-parameters
func (d *Manager) GetList(opts ListOptions) (deployments []*Deployment, err error) {
	res, err := d.client.Get("/deployment", opts)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &deployments)
	return
}

// GetListCount a queries for the number of deployments that fulfill given parameters.
// Takes the same parameters as the Get Deployments method
func (d *Manager) GetListCount(query map[string]string) (count int, err error) {
	res, err := d.client.Get("/deployment/count", query)
	if err != nil {
		return
	}

	resCount := camunda.ResponseCount{}
	err = d.client.Marshal(res, &resCount)
	return resCount.Count, err
}

// Get retrieves a deployment by id, according to the Deployment interface of the engine
func (d *Manager) Get(id string) (deployment Deployment, err error) {
	res, err := d.client.Get("/deployment/"+id, nil)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &deployment)
	return
}

func (d *Manager) writeFields(dc *CreateRequest, w *multipart.Writer) (err error) {
	if err = w.WriteField("deployment-name", dc.DeploymentName); err != nil {
		return
	}

	if dc.EnableDuplicateFiltering {
		if err = w.WriteField("enable-duplicate-filtering", "true"); err != nil {
			return
		}
	}

	if dc.DeployChangedOnly {
		if err = w.WriteField("deploy-changed-only", "true"); err != nil {
			return
		}
	}

	if dc.DeploymentSource != "" {
		if err = w.WriteField("deployment-source", dc.DeploymentSource); err != nil {
			return
		}
	}

	if dc.TenantID != "" {
		if err = w.WriteField("tenant-id", dc.TenantID); err != nil {
			return
		}
	}

	return
}

// Create creates a deployment.
// See more at: https://docs.camunda.org/manual/latest/reference/rest/deployment/post-deployment/
func (d *Manager) Create(dc *CreateRequest) (cr *CreateResponse, err error) {
	cr = &CreateResponse{}

	var data []byte
	body := bytes.NewBuffer(data)
	w := multipart.NewWriter(body)

	d.writeFields(dc, w)

	for key, resource := range dc.Resources {
		var fw io.Writer

		if x, ok := resource.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := resource.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return nil, err
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		}

		if r, ok := resource.(io.Reader); ok {
			if _, err = io.Copy(fw, r); err != nil {
				return nil, err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	res, err := d.client.Post("/deployment/create", nil, body, w.FormDataContentType())
	if err != nil {
		return nil, err
	}

	err = d.client.Marshal(res, cr)

	return cr, err
}

// Redeploy a re-deploys an existing deployment.
// The deployment resources to re-deploy can be restricted by using the properties resourceIds or resourceNames.
// If no deployment resources to re-deploy are passed then all existing resources of the given deployment
// are re-deployed
func (d *Manager) Redeploy(id string, req RedeployRequest) (deployment *CreateResponse, err error) {
	deployment = &CreateResponse{}
	res, err := d.client.Post("/deployment/"+id+"/redeploy", map[string]string{}, &req)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, deployment)
	return
}

// GetResources retrieves all deployment resources of a given deployment
func (d *Manager) GetResources(id string) (resources []*ResourceResponse, err error) {
	res, err := d.client.Get("/deployment/"+id+"/resources", nil)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &resources)
	return
}

// GetResource retrieves a deployment resource by resource id for the given deployment
func (d *Manager) GetResource(id, resourceID string) (resource *ResourceResponse, err error) {
	resource = &ResourceResponse{}
	res, err := d.client.Get("/deployment/"+id+"/resources/"+resourceID, nil)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &resource)
	return
}

// GetResourceBinary retrieves the binary content of a deployment resource for the given deployment by id
func (d *Manager) GetResourceBinary(id, resourceID string) (data []byte, err error) {
	res, err := d.client.Get("/deployment/"+id+"/resources/"+resourceID+"/data", nil)
	if err != nil {
		return
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// Delete deletes a deployment by id
func (d *Manager) Delete(id string, options *DeleteOptions) error {
	_, err := d.client.Delete("/deployment/"+id, options)
	return err
}
