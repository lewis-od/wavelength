package lerna

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/builder"
	"github.com/lewis-od/lambda-build/pkg/executor"
)

type lernaBuilder struct {
	projectName string
	executor    executor.CommandExecutor
}

func NewLerna(commandExecutor executor.CommandExecutor, projectName string) builder.Builder {
	return &lernaBuilder{
		projectName: projectName,
		executor:    commandExecutor,
	}
}

func (l *lernaBuilder) BuildLambda(lambdaName string) (err error) {
	scope := fmt.Sprintf("@%s/%s", l.projectName, lambdaName)
	err = l.executor.Execute([]string{"run", "build", "--scope", scope, "--include-dependencies"})
	return
}
