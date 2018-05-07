package main

import yaml "github.com/sanathkr/go-yaml"

type ParsedTemplate struct {
	Parameters map[string]ParsedParameterSpec `yaml:"Parameters"`
}

type ParsedParameterSpec struct {
	Default *string `yaml:"Default",omitempty`
}

type ParameterSpec struct {
	Name       string
	HasDefault bool
}

func parseTemplate(body []byte) (map[string]ParameterSpec, error) {
	var t ParsedTemplate
	err := yaml.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}
	specs := make(map[string]ParameterSpec)
	for name, parsed := range t.Parameters {
		specs[name] = ParameterSpec{Name: name, HasDefault: parsed.Default != nil}
	}
	return specs, nil
}
