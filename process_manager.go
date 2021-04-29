package camunda

import (
    "io/ioutil"
    "fmt"
)

// ProcessManager a client for ProcessManager
type ProcessManager struct {
    client *Client
}

// GetActivityInstanceStatistics retrieves runtime statistics of a given process definition, grouped by activities.
// These statistics include the number of running activity instances, optionally the number of failed jobs
// and also optionally the number of incidents either grouped by incident types or for a specific incident type.
// Note: This does not include historic data
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/get-activity-statistics/#query-parameters
func (p *ProcessManager) GetActivityInstanceStatistics(by ProcessConfig, query map[string]string) (statistic []*ResActivityInstanceStatistics, err error) {
    res, err := p.client.Get("/process-definition/"+by.Path()+"/statistics", query)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &statistic)
    return
}

// GetDiagram retrieves the diagram of a process definition.
// If the process definition’s deployment contains an image resource with the same file name as the process definition,
// the deployed image will be returned by the Get Diagram endpoint. Example: someProcess.bpmn and someProcess.png.
// Supported file extentions for the image are: svg, png, jpg, and gif
func (p *ProcessManager) GetDiagram(by ProcessConfig) (data []byte, err error) {
    res, err := p.client.Get("/process-definition/"+by.Path()+"/diagram", nil)
    if err != nil {
        return
    }

    defer res.Body.Close()
    return ioutil.ReadAll(res.Body)
}

// GetStartFormVariables Retrieves the start form variables for a process definition
// (only if they are defined via the Generated Task Form approach). The start form variables take form data specified
// on the start event into account. If form fields are defined, the variable types and default values of the form
// fields are taken into account
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/get-form-variables/#query-parameters
func (p *ProcessManager) GetStartFormVariables(by ProcessConfig, filter *FormVariableFilter) (variables map[string]Variable, err error) {
    res, err := p.client.Get("/process-definition/"+by.Path()+"/form-variables", filter)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &variables)
    return
}

// GetListCount requests the number of process definitions that fulfill the query criteria.
// Takes the same filtering parameters as the Get Definitions method
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/get-query-count/#query-parameters
func (p *ProcessManager) GetListCount(query map[string]string) (count int, err error) {
    resCount := ResponseCount{}
    res, err := p.client.Get("/process-definition/count", query)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &resCount)
    return resCount.Count, err
}

// GetList queries for process definitions that fulfill given parameters.
// Parameters may be the properties of process definitions, such as the name, key or version.
// The size of the result set can be retrieved by using the Get Definition Count method
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/get-query/#query-parameters
func (p *ProcessManager) GetList(query map[string]string) (processDefinitions []*ProcessDefinitionResponse, err error) {
    res, err := p.client.Get("/process-definition", query)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &processDefinitions)
    return
}

// GetRenderedStartForm retrieves the rendered form for a process definition.
// This method can be used for getting the HTML rendering of a Generated Task Form
func (p *ProcessManager) GetRenderedStartForm(by ProcessConfig) (htmlForm string, err error) {
    res, err := p.client.Get("/process-definition/"+by.Path()+"/rendered-form", nil)
    if err != nil {
        return
    }

    defer res.Body.Close()
    rawData, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return
    }

    return string(rawData), nil
}

// GetStartFormKey retrieves the key of the start form for a process definition.
// The form key corresponds to the FormData#formKey property in the engine
func (p *ProcessManager) GetStartFormKey(by ProcessConfig) (resp *ResGetStartFormKey, err error) {
    resp = &ResGetStartFormKey{}
    res, err := p.client.Get("/process-definition/"+by.Path()+"/startForm", nil)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &resp)
    return
}

// GetProcessInstanceStatistics retrieves runtime statistics of the process engine, grouped by process definitions.
// These statistics include the number of running process instances, optionally the number of failed jobs and also optionally the number of incidents either grouped by incident types or for a specific incident type.
// Note: This does not include historic data
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/get-statistics/#query-parameters
func (p *ProcessManager) GetProcessInstanceStatistics(query map[string]string) (statistic []*ResInstanceStatistics, err error) {
    res, err := p.client.Get("/process-definition/statistics", query)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &statistic)
    return
}

// GetXML retrieves the BPMN 2.0 XML of a process definition
func (p *ProcessManager) GetXML(by ProcessConfig) (resp *ResBPMNProcessDefinition, err error) {
    resp = &ResBPMNProcessDefinition{}
    res, err := p.client.Get("/process-definition/"+by.Path()+"/xml", nil)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &resp)
    return
}

// Get retrieves a process definition according to the ProcessDefinition interface in the engine
func (p *ProcessManager) Get(by ProcessConfig) (processDefinition *ProcessDefinitionResponse, err error) {
    processDefinition = &ProcessDefinitionResponse{}
    res, err := p.client.Get("/process-definition/"+by.Path(), nil)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, &processDefinition)
    return
}

// StartInstance instantiates a given process definition. Process variables and business key may be supplied
// in the request body
func (p *ProcessManager) StartInstance(config ProcessConfig, req InstanceParams) (pd *ProcessDefinition, err error) {
    pd = &ProcessDefinition{}
    res, err := p.client.Post("/process-definition/"+config.Path()+"/start", nil, &req)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, pd)
    return
}

// SubmitStartForm starts a process instance using a set of process variables and the business key.
// If the start event has Form Field Metadata defined, the process engine will perform backend validation for any form
// fields which have validators defined. See Documentation on Generated Task Forms
func (p *ProcessManager) SubmitStartForm(by ProcessConfig, req ReqSubmitStartForm) (reps *ResSubmitStartForm, err error) {
    reps = &ResSubmitStartForm{}
    res, err := p.client.Post("/process-definition/"+by.Path()+"/submit-form", map[string]string{}, &req)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, reps)
    return
}

// ActivateOrSuspendById activates or suspends a given process definition by id or by latest version
// of process definition key
func (p *ProcessManager) ActivateOrSuspendById(by ProcessConfig, req ReqActivateOrSuspendById) error {
    return p.client.doPutJSON("/process-definition/"+by.Path()+"/suspended", map[string]string{}, &req)
}

// ActivateOrSuspendByKey activates or suspends process definitions with the given process definition key
func (p *ProcessManager) ActivateOrSuspendByKey(req ReqActivateOrSuspendByKey) error {
    return p.client.doPutJSON("/process-definition/suspended", map[string]string{}, &req)
}

// UpdateHistoryTimeToLive updates history time to live for process definition.
// The field is used within History cleanup
func (p *ProcessManager) UpdateHistoryTimeToLive(by ProcessConfig, historyTimeToLive int) error {
    return p.client.doPutJSON("/process-definition/"+by.Path()+"/history-time-to-live", map[string]string{}, &map[string]int{"historyTimeToLive": historyTimeToLive})
}

// Delete deletes a process definition from a deployment by id
// https://docs.camunda.org/manual/latest/reference/rest/process-definition/delete-process-definition/#query-parameters
func (p *ProcessManager) Delete(by ProcessConfig, query map[string]string) error {
    _, err := p.client.Delete("/process-definition/"+by.Path(), query)
    return err
}

// GetDeployedStartForm retrieves the deployed form that can be referenced from a start event. For further information please refer to User Guide
func (p *ProcessManager) GetDeployedStartForm(by ProcessConfig) (htmlForm string, err error) {
    res, err := p.client.Get("/process-definition/"+by.Path()+"/deployed-start-form", nil)
    if err != nil {
        return
    }

    defer res.Body.Close()
    rawData, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return
    }

    return string(rawData), nil
}

// RestartProcessInstance restarts process instances that were canceled or terminated synchronously.
// To execute the restart asynchronously, use the Restart Process Instance Async method
// For more information about the difference between synchronous and asynchronous execution,
// please refer to the related section of the user guide
func (p *ProcessManager) RestartProcessInstance(id string, req RestartInstanceRequest) error {
    _, err := p.client.Post("/process-definition/"+id+"/restart", nil, &req)
    return err
}

// RestartProcessInstanceAsync restarts process instances that were canceled or terminated asynchronously.
// To execute the restart synchronously, use the Restart Process Instance method
// For more information about the difference between synchronous and asynchronous execution,
// please refer to the related section of the user guide
func (p *ProcessManager) RestartProcessInstanceAsync(id string, req RestartInstanceRequest) (resp *ResBatch, err error) {
    resp = &ResBatch{}
    res, err := p.client.Post("/process-definition/"+id+"/restart-async", nil, &req)
    if err != nil {
        return
    }

    err = p.client.Marshal(res, resp)
    return
}

// ListInstances Queries for process instances that fulfill given parameters. Parameters may be static as well as
// dynamic runtime properties of process instances. The size of the result set can be retrieved by using the Get
// Instance Count method.
func (p *ProcessManager) ListInstances(q ProcessInstanceQuery) ([]*ProcessInstance, error) {
    var pi []*ProcessInstance

    res, err := p.client.Get("/process-instance", q)
    if err != nil {
        return nil, fmt.Errorf("cannot invoke client: %w", err)
    }
    defer res.Body.Close()

    err = p.client.Marshal(res, &pi)
    if err != nil {
        return nil, fmt.Errorf("cannot unmarshal response: %w", err)
    }

    return pi, nil
}

// GetInstanceVars Retrieves all variables of a given process instance by id.
func (p *ProcessManager) GetInstanceVars(instanceId string) (map[string]Variable, error) {
    res, err := p.client.Get(fmt.Sprintf("/process-instances/%s/variables", instanceId), nil)
    if err != nil {
        return nil, fmt.Errorf("cannot invoke client: %w", err)
    }
    defer res.Body.Close()

    vars := make(map[string]Variable)

    err = p.client.Marshal(res, &vars)
    if err != nil {
        return nil, fmt.Errorf("cannot marshal variables: %w", err)
    }

    return vars, nil
}
