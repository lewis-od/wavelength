package lerna

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/builder"
	"github.com/lewis-od/lambda-build/pkg/executor"
)

type lernaBuilder struct {
	Project string
	executor executor.CommandExecutor
}

func NewLerna(commandExecutor executor.CommandExecutor, projectName string) builder.Builder {
	return &lernaBuilder{
		Project:  projectName,
		executor: commandExecutor,
	}
}

func (l *lernaBuilder) BuildLambda(lambdaName string) (err error) {
	scope := fmt.Sprintf("@%s/%s", l.Project, lambdaName)
	err = l.executor.Execute([]string{"run", "build", "--scope", scope, "--include-dependencies"})
	return
}

func (l *lernaBuilder) BuildLambdas(lambdaNames []string) error {
	for _, lambdaName := range lambdaNames {
		fmt.Printf("🔨 Building %s lambda...\n", lambdaName)
		err := l.BuildLambda(lambdaName)
		if err != nil {
			return fmt.Errorf("❌ Error building %s\n%s\n", lambdaName, err)
		}
	}
	return nil
}
