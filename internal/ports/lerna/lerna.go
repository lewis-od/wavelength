package lerna

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/executor"
)

type lernaBuilder struct {
	projectName *string
	executor    executor.CommandExecutor
}

func NewLerna(commandExecutor executor.CommandExecutor, projectName *string) builder.Builder {
	return &lernaBuilder{
		projectName: projectName,
		executor:    commandExecutor,
	}
}

func (l *lernaBuilder) BuildLambda(lambdaName string) (output []byte, err error) {
	scope := fmt.Sprintf("@%s/%s", *l.projectName, lambdaName)
	cmdContext := &executor.CommandContext{
		Directory: ".",
	}
	output, err = l.executor.ExecuteAndCapture([]string{"run", "build", "--scope", scope, "--include-dependencies"}, cmdContext)
	return
}
