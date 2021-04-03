package terraform

import (
	"encoding/json"
	"github.com/lewis-od/lambda-build/pkg/executor"
)

type Terraform interface {
	Output(directory string) (map[string]Output, error)
}

type TerraformExecutor struct {
	executor executor.CommandExecutor
}

func NewTerraform(executor executor.CommandExecutor) Terraform {
	return &TerraformExecutor{
		executor: executor,
	}
}

type Output struct {
	Sensitive bool `json:"sensitive"`
	Type string `json:"type"`
	Value string `json:"value"`
}

func (tf *TerraformExecutor) Output(directory string) (map[string]Output, error) {
	context := &executor.CommandContext{
		Directory: directory,
	}
	commandOutput, err := tf.executor.ExecuteAndCapture([]string{"output", "-json"}, context)
	if err != nil {
		return nil, err
	}

	var outputs map[string]Output
	err = json.Unmarshal(commandOutput, &outputs)
	return outputs, err
}
