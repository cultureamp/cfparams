package main

import (
	"encoding/json"

	yaml "github.com/sanathkr/go-yaml"
)

type jsonItem struct {
	Key   string
	Value string
}

func getJsonForInputTags(input *Input) ([]byte, error) {
	data := map[string]string{}
	yaml.Unmarshal(input.TagsBody, &data)
	items := []jsonItem{}
	for key, value := range data {
		items = append(items, jsonItem{Key: key, Value: value})
	}
	for key, value := range input.Parameters {
		items = append(items, jsonItem{Key: key, Value: value})
	}
	return json.MarshalIndent(items, "", "  ")
}
