package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Input struct {
	TemplateBody   []byte
	AcceptDefaults bool
	NoPrevious     bool
	ParametersCLI  []string
	ParametersYAML []byte
	Parameters     map[string]string
}

type ParameterItem struct {
	ParameterKey     string `json:"ParameterKey,omitempty"`
	ParameterValue   string `json:"ParameterValue,omitempty"`
	UsePreviousValue bool   `json:"UsePreviousValue,omitempty"`
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

func main() {
	input := &Input{}
	var tplFile, paramFile string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [flags] [Key=value ...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Produces JSON suitable for `aws cloudformation` CLI.\n\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&tplFile, "template", "", "CloudFormation YAML template path")
	flag.StringVar(&paramFile, "parameters", "", "Parameters YAML file")
	flag.BoolVar(&input.AcceptDefaults, "accept-defaults", false, "Accept defaults from CloudFormation template, omit from JSON")
	flag.BoolVar(&input.NoPrevious, "no-previous", false, "Disable UsePreviousValue, fail if a parameter has no default and is not specified")
	flag.Parse()

	if tplFile != "" {
		data, err := ioutil.ReadFile(tplFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read CloudFormation template: %s\n", tplFile)
			os.Exit(1)
		}
		input.TemplateBody = data
	} else {
		fmt.Fprintf(os.Stderr, "CloudFormation template required, e.g: --template=cfn.yaml\n")
		os.Exit(1)
	}

	if paramFile != "" {
		data, err := ioutil.ReadFile(paramFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read parameters file: %s\n", paramFile)
			os.Exit(1)
		}
		input.ParametersYAML = data
	}

	// remaining positional args
	input.ParametersCLI = flag.Args()

	j, err := getJsonForInput(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(j)
	os.Stdout.Write([]byte("\n"))
}

func getJsonForInput(input *Input) ([]byte, error) {
	if err := parseParameters(input); err != nil {
		return nil, err
	}

	var t ParsedTemplate
	t.Parameters = make(map[string]ParsedParameterSpec)
	err := yaml.Unmarshal(input.TemplateBody, &t)
	if err != nil {
		return nil, err
	}

	specs := make(map[string]ParameterSpec)
	for name, parsed := range t.Parameters {
		specs[name] = ParameterSpec{Name: name, HasDefault: parsed.Default != ""}
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
