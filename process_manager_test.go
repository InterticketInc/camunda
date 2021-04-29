package camunda

import (
    "testing"
    "encoding/json"
)

func createTestClient() *Client {
    c := NewClient(&ClientOptions{EndpointUrl: DefaultEndpointUrl})
    //dm := deploy.NewManager(c)
    //dm.Create(&deploy.CreateRequest{
    //    DeploymentName:           "",
    //    EnableDuplicateFiltering: false,
    //    DeployChangedOnly:        false,
    //    DeploymentSource:         "",
    //    TenantID:                 "",
    //    Resources:                nil,
    //})

    return c
}

func TestProcessManager_GetInstanceVars(t *testing.T) {
    ret := `{"aVariableKey": {
        "value" : {"prop1" : "a", "prop2" : "b"},
        "type" : "Object",
        "valueInfo" : {
          "objectTypeName": "com.example.MyObject",
          "serializationDataFormat": "application/xml"
        }
      }}`

    m := make(map[string]Variable)
    err := json.Unmarshal([]byte(ret), &m)
    if err != nil {
        t.Fatalf("cannot unmarshal %s", err.Error())
    }

    t.Logf("macika: %+v", m)
}

func TestProcessManager_ListInstances(t *testing.T) {
    pm := createTestClient().ProcessManager()

    //_, _ = pm.StartInstance(ProcessConfig{Key: "test-key", TenantId: "test"}, InstanceParams{
    //    BusinessKey: "test-instance-key",
    //})

    instances, err := pm.ListInstances(ProcessInstanceQuery{
        BusinessKey: "test-instance-key",
    })

    if err != nil {
        t.Fatalf("cannot list instances: %s", err.Error())
    }

    t.Logf("Instance list size: %d", len(instances))

    for _, instance := range instances {
        t.Logf("%s", instance.BusinessKey)
    }
    t.Logf("Instances: %+v", instances)
}
