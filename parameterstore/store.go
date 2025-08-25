package parameterstore

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var fakes map[string]string

func Fake(params map[string]string) {
	fakes = params
}

func Get(name string) (string, error) {
	if fakes == nil {
		return getReal(name)
	} else {
		return getFake(name)
	}
}

func getFake(name string) (string, error) {
	if val, ok := fakes[name]; ok {
		return val, nil
	} else {
		return "", errors.New("No fake parameter for " + name)
	}
}

func getReal(name string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}
	client := ssm.NewFromConfig(cfg)
	output, err := client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, nil
}
