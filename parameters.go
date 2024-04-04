package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/cultureamp/cfparams/parameterstore"
	yaml "github.com/sanathkr/go-yaml"
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

type parameterStoreUnmarshaler struct{}

func (t *parameterStoreUnmarshaler) UnmarshalYAMLTag(tag string, fieldValue reflect.Value) reflect.Value {
	name := fieldValue.String()
	log.New(os.Stderr, "", log.LstdFlags).Printf("ParameterStore: GetParameter(%#v)\n", name)
	value, err := parameterstore.Get(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		value = "" // crash instead?
	}
	return reflect.ValueOf(value)
}

func getJsonForInputParams(input *Input) ([]byte, error) {
	if err := parseParameters(input); err != nil {
		return nil, err
	}

	specs, err := parseTemplate(input.TemplateBody)
	if err != nil {
		return nil, err
	}

	if err := validateParameters(input.Parameters, specs); err != nil {
		return nil, err
	}

	items := []ParameterItem{}
	missingNames := []string{}
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
			missingNames = append(missingNames, spec.Name)
		}
	}

	if len(missingNames) > 0 {
		return nil, fmt.Errorf("missing parameters: %s", strings.Join(missingNames, ", "))
	}

	return json.MarshalIndent(items, "", "  ")
}

func parseParameters(input *Input) error {
	input.Parameters = make(map[string]string)

	yaml.RegisterTagUnmarshaler("!ParameterStore", &parameterStoreUnmarshaler{})

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
	for name := range params {
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
