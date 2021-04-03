package system

import (
	"github.com/lewis-od/lambda-build/pkg/executor"
	"os"
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

func (se *Executor) Execute(args []string) (err error) {
	cmd := exec.Command(se.command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func (se *Executor) ExecuteAndCapture(args []string, context *executor.CommandContext) ([]byte, error) {
	cmd := exec.Command(se.command, args...)
	if context.Directory != "" {
		cmd.Dir = context.Directory
	}
	output, err := cmd.Output()
	return output, err
}
