package main

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "github.com/sanathkr/go-yaml"
)

type jsonItem struct {
	Key   string
	Value string
}

func getJsonForInputTags(input *Input) ([]byte, error) {
	data := map[string]string{}

	// Tags from YAML file
	err := yaml.Unmarshal(input.TagsBody, &data)
	if err != nil {
		return nil, err
	}

	// Tags from CLI
	for _, kv := range input.ParametersCLI {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("expected key=value, got %s", pair)
		}
		data[pair[0]] = pair[1]
	}

	items := []jsonItem{}
	for key, value := range data {
		items = append(items, jsonItem{Key: key, Value: value})
	}
	for key, value := range input.Parameters {
		items = append(items, jsonItem{Key: key, Value: value})
	}
	return json.MarshalIndent(items, "", "  ")
}
