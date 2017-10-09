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
	input := Input{TemplateBody: []byte(cfnTemplate)}
	expected := []ParameterItem{
		{ParameterKey: "Greeting", UsePreviousValue: true},
		{ParameterKey: "Recipient", UsePreviousValue: true},
		{ParameterKey: "ImageRepo", UsePreviousValue: true},
		{ParameterKey: "ImageTag", UsePreviousValue: true},
		{ParameterKey: "Cluster", UsePreviousValue: true},
	}

	json, err := parametersJson(input)
	require.NoError(t, err)
	actual, err := itemsFromJson(json)
	require.NoError(t, err)

	assert.Equal(t, len(expected), len(actual))
	for _, item := range actual {
		assert.Contains(t, expected, item)
	}
}

func itemsFromJson(j []byte) ([]ParameterItem, error) {
	var items []ParameterItem
	err := json.Unmarshal(j, &items)
	return items, err
}
