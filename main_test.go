package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cfnTemplate = []byte(`
Parameters:
  Greeting:
    Type: String
    Description: greeting message to send
    Default: Hello
  Recipient:
    Description: name of the greeting recipient
    Type: String
  ImageRepo:
    Type: String
    Description: repository of Docker image to run
    Default: "123.dkr.ecr.us-east-1.amazonaws.com/greeting"
  ImageTag:
    Type: String
    Description: tag of Docker image to run
    Default: latest
  Cluster:
    Description: ECS cluster ID to run service on
    Type: String
`)

func TestUsePreviousAll(t *testing.T) {
	input := &Input{TemplateBody: []byte(cfnTemplate)}
	actual := mustGetParameterItems(t, input)
	expected := []ParameterItem{
		{ParameterKey: "Greeting", UsePreviousValue: true},
		{ParameterKey: "Recipient", UsePreviousValue: true},
		{ParameterKey: "ImageRepo", UsePreviousValue: true},
		{ParameterKey: "ImageTag", UsePreviousValue: true},
		{ParameterKey: "Cluster", UsePreviousValue: true},
	}
	assert.Equal(t, len(expected), len(actual))
	for _, item := range actual {
		assert.Contains(t, expected, item)
	}
}

func TestLaunchScenarioCLI(t *testing.T) {
	input := &Input{
		TemplateBody:   cfnTemplate,
		AcceptDefaults: true,
		NoPrevious:     true,
		ParametersCLI: []string{
			"Recipient=world",
			"ImageTag=v1",
			"Cluster=nanoservices",
		},
	}
	actual := mustGetParameterItems(t, input)
	expected := []ParameterItem{
		{ParameterKey: "Recipient", ParameterValue: "world"},
		{ParameterKey: "ImageTag", ParameterValue: "v1"},
		{ParameterKey: "Cluster", ParameterValue: "nanoservices"},
	}
	assert.Equal(t, len(expected), len(actual))
	for _, item := range expected {
		assert.Contains(t, actual, item)
	}
}

func TestLaunchScenarioFile(t *testing.T) {
	input := &Input{
		TemplateBody:   cfnTemplate,
		AcceptDefaults: true,
		NoPrevious:     true,
		ParametersCLI:  []string{"ImageTag=v1"},
		ParametersYAML: []byte("---\nRecipient: world\nCluster: nanoservices\n"),
	}
	actual := mustGetParameterItems(t, input)
	expected := []ParameterItem{
		{ParameterKey: "Recipient", ParameterValue: "world"},
		{ParameterKey: "ImageTag", ParameterValue: "v1"},
		{ParameterKey: "Cluster", ParameterValue: "nanoservices"},
	}
	assert.Equal(t, len(expected), len(actual))
	for _, item := range expected {
		assert.Contains(t, actual, item)
	}
}

func TestDeployScenario(t *testing.T) {
	input := &Input{
		TemplateBody:  cfnTemplate,
		ParametersCLI: []string{"ImageTag=v2"},
	}
	actual := mustGetParameterItems(t, input)
	expected := []ParameterItem{
		{ParameterKey: "Greeting", UsePreviousValue: true},
		{ParameterKey: "Recipient", UsePreviousValue: true},
		{ParameterKey: "ImageRepo", UsePreviousValue: true},
		{ParameterKey: "ImageTag", ParameterValue: "v2"},
		{ParameterKey: "Cluster", UsePreviousValue: true},
	}
	assert.Equal(t, len(expected), len(actual))
	for _, item := range expected {
		assert.Contains(t, actual, item)
	}
}

func mustGetParameterItems(t *testing.T, input *Input) []ParameterItem {
	j, err := getJsonForInput(input)
	require.NoError(t, err)
	var items []ParameterItem
	err = json.Unmarshal(j, &items)
	require.NoError(t, err)
	return items
}
