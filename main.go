package main

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Input struct {
	TemplateBody   []byte
	AcceptDefaults bool
	NoPrevious     bool
	ParametersCLI  []string
	Parameters     map[string]string
}

type ParameterItem struct {
	ParameterKey     string `json:"ParameterKey,omitempty"`
	ParameterValue   string `json:"ParameterValue,omitempty"`
	UsePreviousValue bool   `json:UsePreviousValue,omitempty"`
}

type ParsedParameterSpec struct {
	Type        string `yaml:"Type"`
	Description string `yaml:"Description"`
	Default     string `yaml:"Default"`
}

type ParsedTemplate struct {
	Parameters map[string]ParsedParameterSpec `yaml:"Parameters"`
}

type ParameterSpec struct {
	Name       string
	HasDefault bool
}

func parametersJson(input Input) ([]byte, error) {
	input.Parameters = make(map[string]string)
	for _, kv := range input.ParametersCLI {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("expected Key=value, got %s", pair)
		}
		input.Parameters[pair[0]] = pair[1]
	}

	var t ParsedTemplate
	t.Parameters = make(map[string]ParsedParameterSpec)
	err := yaml.Unmarshal(input.TemplateBody, &t)
	if err != nil {
		return nil, err
	}

	specs := []ParameterSpec{}
	for name, parsed := range t.Parameters {
		specs = append(specs, ParameterSpec{Name: name, HasDefault: parsed.Default != ""})
	}

	items := []ParameterItem{}
	for _, spec := range specs {
		if value, ok := input.Parameters[spec.Name]; ok {
			// specified in parameters
			items = append(items, ParameterItem{
				ParameterKey:   spec.Name,
				ParameterValue: value,
			})
		} else if input.AcceptDefaults && spec.HasDefault {
			// has default; do not override
			continue
		} else if !input.NoPrevious {
			// use previous value
			items = append(items, ParameterItem{
				ParameterKey:     spec.Name,
				UsePreviousValue: true,
			})
		} else {
			return nil, fmt.Errorf("no parameter found for %s", spec.Name)
		}
	}

	return json.Marshal(items)
}
