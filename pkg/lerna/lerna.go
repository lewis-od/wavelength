package lerna

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/executor"
)

type Lerna interface {
	BuildLambda(lambdaName string) (err error)
}

type lernaExecutor struct {
	Project string
	executor executor.CommandExecutor
}

func NewLerna(commandExecutor executor.CommandExecutor, projectName string) Lerna {
	return &lernaExecutor{
		Project:  projectName,
		executor: commandExecutor,
	}
}

func (l *lernaExecutor) BuildLambda(lambdaName string) (err error) {
	scope := fmt.Sprintf("@%s/%s", l.Project, lambdaName)
	err = l.executor.Execute([]string{"run", "build", "--scope", scope, "--include-dependencies"})
	return
}
