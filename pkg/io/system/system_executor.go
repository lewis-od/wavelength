package system

import (
	"github.com/lewis-od/lambda-build/pkg/executor"
	"os/exec"
)

type Executor struct {
	command string
}

func NewExecutor(command string) *Executor {
	return &Executor{
		command: command,
	}
}

func (se *Executor) ExecuteWithContext(args []string, context *executor.CommandContext) ([]byte, error) {
	cmd := exec.Command(se.command, args...)
	if context.Directory != "" {
		cmd.Dir = context.Directory
	}
	output, err := cmd.Output()
	return output, err
}
