package camunda

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MessageManager struct {
	client *Client
}

// NewMessageManager initializes the
func NewMessageManager(client *Client) *MessageManager {
	return &MessageManager{client: client}
}

// SendMessage Correlates a message to the process engine to either trigger a message start event or an intermediate
// message catching event. Internally this maps to the engine’s message correlation builder methods
// MessageCorrelationBuilder#correlateWithResult() and MessageCorrelationBuilder#correlateAllWithResult().
// For more information about the correlation behavior, see the Message Events section of the BPMN 2.0
// Implementation Reference.
func (mm *MessageManager) SendMessage(request *MessageRequest) (*SendMessageResponse, error) {
	res, err := mm.client.Post("/message", nil, request)
	if err != nil {
		return nil, fmt.Errorf("cannot send message: %w", err)
	}
	defer res.Body.Close()

	// No content
	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var smr SendMessageResponse
	err = json.NewDecoder(res.Body).Decode(&smr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode request: %w", err)
	}

	return &smr, nil
}

type MessageRequest struct {
	// MessageName	The name of the message to deliver.
	MessageName string `json:"messageName"`
	// BusinessKey	Used for correlation of process instances that wait for incoming messages.
	// Will only correlate to executions that belong to a process instance with the provided business key.
	// If the message triggers a start event, the business key is set in the process instance.
	BusinessKey string `json:"businessKey"`
	// TenantID	Used to correlate the message for a tenant with the given id. Will only correlate to executions and
	// process definitions which belong to the tenant. Must not be supplied in conjunction with a withoutTenantId.
	TenantID string `json:"tenantId,omitempty"`
	// WithoutTenantID A Boolean value that indicates whether the message should only be correlated to executions and
	// process definitions which belong to no tenant or not.
	// Value may only be true, as false is the default behavior.Must not be supplied in conjunction with a tenantId.
	WithoutTenantID string `json:"withoutTenantId,omitempty"`
	// ProcessInstanceID    Used to correlate the message to the process instance with the given id.
	ProcessInstanceID string `json:"processInstanceId,omitempty"`
	// CorrelationKeys Used for correlation of process instances that wait for incoming messages.
	// Has to be a JSON object containing key-value pairs that are matched against process instance variables during
	// correlation.Each key is a variable name and each value a JSON variable value object with the following
	// properties.
	CorrelationKeys Variables `json:"correlationKeys,omitempty"`
	// LocalCorrelationKeys Local variables used for correlation of executions (process instances) that wait for
	// incoming messages. Has to be a JSON object containing key-value pairs that are matched against local variables
	// during correlation.Each key is a variable name and each value a JSON variable value object with the following
	// properties.
	LocalCorrelationKeys Variables `json:"localCorrelationKeys,omitempty"`
	// ProcessVariables A map of variables that is injected into the triggered execution or process instance after
	// the message has been delivered. Each key is a variable name and each value a JSON variable value object with the
	// following properties.
	ProcessVariables Variables `json:"processVariables,omitempty"`
	// ProcessVariablesLocal    A map of local variables that is injected into the triggered execution or process
	// instance after the message has been delivered.Each key is a variable name and each value a JSON variable value
	// object with the following properties.
	ProcessVariablesLocal Variables `json:"processVariablesLocal"`
	// All A Boolean value that indicates whether the message should be correlated to exactly one entity or multiple
	// entities. If the value is set to false, the message will be correlated to exactly one entity (execution or
	// process definition).If the value is set to true, the message will be correlated to multiple executions and a
	// process definition that can be instantiated by this message in one go.
	All bool `json:"all,omitempty"`
	// ResultEnabled A Boolean value that indicates whether the result of the correlation should be returned or not.
	// If this property is set to true, there will be returned a list of message correlation result objects. Depending
	// on the all property, there will be either one ore more returned results in the list.
	// The default value is false, which means no result will be returned.
	ResultEnabled bool `json:"resultEnabled"`
	// VariablesInResultEnabled    A Boolean value that indicates whether the result of the correlation should contain
	// process variables or not.The parameter resultEnabled should be set to true in order to use this it.
	// The default value is false, which means the variables will not be returned.
	VariablesInResultEnabled bool `json:"variablesInResultEnabled"`
}

type SendMessageResponse struct {
	// ResultType String Indicates if the message was correlated to a message start event or an intermediate message
	// catching event. In the first case, the resultType is ProcessDefinition and otherwise Execution.
	ResultType string `json:"resultType,omitempty"`
	// ProcessInstance	Object	This property only has a value if the resultType is set to ProcessDefinition.
	// The processInstance with the properties as described in the get single instance method.
	ProcessInstance string `json:"processInstance,omitempty"`
	// Execution Object	This property only has a value if the resultType is set to Execution. The execution
	// with the properties as described in the get single execution method.
	Execution interface{} `json:"execution,omitempty"`
	// Variables This property is returned if the variablesInResultEnabled is set to true.
	// Contains a list of the process variables.
	Variables Variables
}
