package terraform

import (
	"encoding/json"
	"github.com/lewis-od/lambda-build/pkg/executor"
)

type Terraform struct {
	executor executor.CommandExecutor
}

func NewTerraform(executor executor.CommandExecutor) *Terraform {
	return &Terraform{
		executor: executor,
	}
}

type Output struct {
	Sensitive bool `json:"sensitive"`
	Type string `json:"type"`
	Value string `json:"value"`
}

func (tf *Terraform) output(directory string) (map[string]Output, error) {
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
