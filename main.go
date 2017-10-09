package main

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

type Input struct {
	TemplateBody []byte
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
		items = append(items, ParameterItem{
			ParameterKey:     spec.Name,
			UsePreviousValue: true,
		})
	}

	return json.Marshal(items)
}
