package camunda

import (
	"fmt"
	"strings"
)

type Variables map[string]Variable

func (v Variables) String(name string) (string, error) {
	if strings.ToLower(v[name].Type) != "string" {
		return "", fmt.Errorf("cannot convert value type %s to string", v[name].Type)
	}

	return v[name].Value.(string), nil
}

func (v Variables) Int(name string) (int, error) {
	if strings.ToLower(v[name].Type) != "Integer" {
		return -1, fmt.Errorf("cannot convert value type %s to Integer", v[name].Type)
	}

	return v[name].Value.(int), nil
}
