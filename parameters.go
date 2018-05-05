package main

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type ParameterItem struct {
	ParameterKey     string
	ParameterValue   string
	UsePreviousValue bool
}

// JSON representation for parameter with ParameterValue
type ParameterItemWithValue struct {
	ParameterKey   string `json:"ParameterKey"`
	ParameterValue string `json:"ParameterValue"`
}

// JSON representation for parameter with UsePreviousValue
type ParameterItemUsePrevious struct {
	ParameterKey     string `json:"ParameterKey"`
	UsePreviousValue bool   `json:"UsePreviousValue"`
}

func parseParameters(input *Input) error {
	input.Parameters = make(map[string]string)

	// Parameters from YAML file
	err := yaml.Unmarshal(input.ParametersYAML, input.Parameters)
	if err != nil {
		return err
	}

	// Parameters from CLI
	for _, kv := range input.ParametersCLI {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			return fmt.Errorf("expected Key=value, got %s", pair)
		}
		input.Parameters[pair[0]] = pair[1]
	}

	return nil
}

func validateParameters(params map[string]string, specs map[string]ParameterSpec) error {
	unexpected := []string{}
	for name, _ := range params {
		if _, ok := specs[name]; !ok {
			unexpected = append(unexpected, name)
		}
	}
	if len(unexpected) > 0 {
		return fmt.Errorf("specified parameters not in template: %s", strings.Join(unexpected, ", "))
	}
	return nil
}

func (p ParameterItem) MarshalJSON() ([]byte, error) {
	if p.UsePreviousValue {
		return json.Marshal(ParameterItemUsePrevious{
			ParameterKey:     p.ParameterKey,
			UsePreviousValue: true,
		})
	} else {
		return json.Marshal(ParameterItemWithValue{
			ParameterKey:   p.ParameterKey,
			ParameterValue: p.ParameterValue,
		})
	}
}
