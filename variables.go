package camunda

import (
    "encoding/json"
    "fmt"
    "strings"

    "github.com/mitchellh/mapstructure"
    "errors"
)

var (
    ErrContextVarNotFound = errors.New("context variable not found")
)

type Variables map[string]*Variable

func (v Variables) String(name string) (string, error) {
    if _, ok := v[name]; !ok {
        return "", fmt.Errorf("variable '%s' not found", name)
    }

    if strings.ToLower(v[name].Type) != "string" {
        return "", fmt.Errorf("cannot convert value type %s to string", v[name].Type)
    }

    return v[name].Value.(string), nil
}

func (v Variables) Int(name string) (int, error) {
    if _, ok := v[name]; !ok {
        return -1, fmt.Errorf("variable '%s' not found", name)
    }

    t := strings.ToLower(v[name].Type)
    if t != "Integer" && t != "Long" {
        return -1, fmt.Errorf("cannot convert value type %s to Integer", v[name].Type)
    }

    return v[name].Value.(int), nil
}

func (v Variables) JSON(name string) ([]byte, error) {
    if _, ok := v[name]; !ok {
        return nil, fmt.Errorf("variable '%s' not found", name)
    }

    if !strings.EqualFold(v[name].Type, "Json") {
        return nil, fmt.Errorf("cannot convert value type %s to Json", v[name].Type)
    }

    return []byte(v[name].Value.(string)), nil
}

// MarshalField marshals a field to the destination interface
// nolint:wrapcheck
func (v Variables) MarshalField(name string, field interface{}) error {
    decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
        TagName: "json",
        Result:  field,
    })

    if err != nil {
        return fmt.Errorf("cannot create decoder: %w", err)
    }

    val, ok := v[name]
    if !ok {
        return ErrContextVarNotFound
    }

    switch val.Value.(type) {
    case string:
        return json.Unmarshal([]byte(val.Value.(string)), field)
    default:
        return decoder.Decode(val.Value)
    }
}

// Map mapping values to original map format without type definitions and valueInfo fields
func (v Variables) Map() map[string]interface{} {
    m := make(map[string]interface{}, 0)
    for k, val := range v {
        m[k] = val.Value
    }

    return m
}

func (v Variables) AddString(key string, value string) {
    v[key] = &Variable{
        Value: value,
        Type:  "String",
    }
}

func (v Variables) AddJSON(key string, val interface{}) {
    bb, _ := json.Marshal(val)
    v.AddJSONBytes(key, bb)
}

func (v Variables) AddJSONBytes(key string, bb []byte) {
    v[key] = &Variable{
        Value: string(bb),
        Type:  "Json",
    }
}

func (v Variables) AddList(key string, values interface{}) {
    jsonString, _ := json.Marshal(values)

    v[key] = &Variable{
        Value: string(jsonString),
        Type:  "Object",
        ValueInfo: &ValueInfo{
            ObjectTypeName:          "java.util.ArrayList",
            SerializationDataFormat: "application/json",
        },
    }
}

func (v Variables) AddObject(key string, value interface{}) {
    jsonString, _ := json.Marshal(value)

    v[key] = &Variable{
        Value: string(jsonString),
        Type:  "Object",
        ValueInfo: &ValueInfo{
            ObjectTypeName:          "java.util.LinkedHashMap",
            SerializationDataFormat: "application/json",
        },
    }
}

func CreateVariables(v map[string]interface{}) Variables {
    vars := Variables{}
    for k, val := range v {
        tmp := Variable{
            Value: val,
            Type:  "Json",
        }

        switch val.(type) {
        case float64, float32, int, int32, int64, int8, int16:
            tmp.Type = "Long"
        case string:
            tmp.Type = "String"
        case []byte: // already marshalled object
            tmp.Value = string(val.([]byte))
        default:
            // Marshal object to JSON
            bb, _ := json.Marshal(val)
            tmp.Value = string(bb)
        }

        vars[k] = &tmp
    }

    return vars
}
