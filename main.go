package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// see Makefile
var Version = "dev"

type Input struct {
	TemplateBody   []byte
	AcceptDefaults bool
	NoPrevious     bool
	ParametersCLI  []string
	ParametersYAML []byte
	Parameters     map[string]string
}

func main() {
	input := &Input{}
	var tplFile, paramFile string

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "cfparams %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [Key=value ...]\n\n", os.Args[0])
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
	return getJsonForInputParams(input) // parameters.go
}
