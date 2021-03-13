package camunda

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// DeployManager deployment manager instance. You can instantiate the instance from the camunda.Client instance
type DeployManager struct {
	client *Client
}

// GetList a queries for deployments that fulfill given parameters. Parameters may be the properties of deployments,
// such as the id or name or a range of the deployment time. The size of the result set can be retrieved by using
// the Get Deployment count method.
// Query parameters described in the documentation:
// https://docs.camunda.org/manual/latest/reference/rest/deployment/get-query/#query-parameters
func (d *DeployManager) GetList(query map[string]string) (deployments []*ResDeployment, err error) {
	res, err := d.client.Get("/deployment", query)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &deployments)
	return
}

// GetListCount a queries for the number of deployments that fulfill given parameters.
// Takes the same parameters as the Get Deployments method
func (d *DeployManager) GetListCount(query map[string]string) (count int, err error) {
	res, err := d.client.Get("/deployment/count", query)
	if err != nil {
		return
	}

	resCount := ResponseCount{}
	err = d.client.Marshal(res, &resCount)
	return resCount.Count, err
}

// Get retrieves a deployment by id, according to the Deployment interface of the engine
func (d *DeployManager) Get(id string) (deployment ResDeployment, err error) {
	res, err := d.client.Get("/deployment/"+id, map[string]string{})
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &deployment)
	return
}

// Create creates a deployment
func (d *DeployManager) Create(dc *ReqDeploymentCreate) (deployment *ResDeploymentCreate, err error) {
	deployment = &ResDeploymentCreate{}
	var data []byte
	body := bytes.NewBuffer(data)
	w := multipart.NewWriter(body)

	if err = w.WriteField("deployment-name", dc.DeploymentName); err != nil {
		return nil, err
	}

	if dc.EnableDuplicateFiltering != nil {
		if err = w.WriteField("enable-duplicate-filtering", strconv.FormatBool(*dc.EnableDuplicateFiltering)); err != nil {
			return nil, err
		}
	}

	if dc.DeployChangedOnly != nil {
		if err = w.WriteField("deploy-changed-only", strconv.FormatBool(*dc.DeployChangedOnly)); err != nil {
			return nil, err
		}
	}

	if dc.DeploymentSource != nil {
		if err = w.WriteField("deployment-source", *dc.DeploymentSource); err != nil {
			return nil, err
		}
	}

	if dc.TenantID != nil {
		if err = w.WriteField("tenant-id", *dc.TenantID); err != nil {
			return nil, err
		}
	}

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

	res, err := d.client.do(http.MethodPost, "/deployment/create", map[string]string{}, body, w.FormDataContentType())
	if err != nil {
		return nil, err
	}

	err = d.client.Marshal(res, deployment)

	return deployment, err
}

// Redeploy a re-deploys an existing deployment.
// The deployment resources to re-deploy can be restricted by using the properties resourceIds or resourceNames.
// If no deployment resources to re-deploy are passed then all existing resources of the given deployment
// are re-deployed
func (d *DeployManager) Redeploy(id string, req ReqRedeploy) (deployment *ResDeploymentCreate, err error) {
	deployment = &ResDeploymentCreate{}
	res, err := d.client.post("/deployment/"+id+"/redeploy", map[string]string{}, &req)
	if err != nil {
		return
	}

	err = d.client.Marshal(res, deployment)
	return
}

// GetResources retrieves all deployment resources of a given deployment
func (d *DeployManager) GetResources(id string) (resources []*ResDeploymentResource, err error) {
	res, err := d.client.Get("/deployment/"+id+"/resources", map[string]string{})
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &resources)
	return
}

// GetResource retrieves a deployment resource by resource id for the given deployment
func (d *DeployManager) GetResource(id, resourceID string) (resource *ResDeploymentResource, err error) {
	resource = &ResDeploymentResource{}
	res, err := d.client.Get("/deployment/"+id+"/resources/"+resourceID, map[string]string{})
	if err != nil {
		return
	}

	err = d.client.Marshal(res, &resource)
	return
}

// GetResourceBinary retrieves the binary content of a deployment resource for the given deployment by id
func (d *DeployManager) GetResourceBinary(id, resourceID string) (data []byte, err error) {
	res, err := d.client.Get("/deployment/"+id+"/resources/"+resourceID+"/data", map[string]string{})
	if err != nil {
		return
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// Delete deletes a deployment by id
func (d *DeployManager) Delete(id string, query map[string]string) error {
	_, err := d.client.delete("/deployment/"+id, query)
	return err
}
