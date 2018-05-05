package parameterstore

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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
	client := ssm.New(session.Must(session.NewSession(&aws.Config{})))
	output, err := client.GetParameter(&ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return *output.Parameter.Value, nil
}
