package system

import (
	"github.com/lewis-od/lambda-build/pkg/executor"
	"os"
	"os/exec"
)

type sysExecutor struct {
	command string
}

func NewExecutor(command string) executor.CommandExecutor {
	return &sysExecutor{
		command: command,
	}
}

func (se *sysExecutor) Execute(args []string) (err error) {
	cmd := exec.Command(se.command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func (se *sysExecutor) ExecuteAndCapture(args []string, context *executor.CommandContext) ([]byte, error) {
	cmd := exec.Command(se.command, args...)
	if context.Directory != "" {
		cmd.Dir = context.Directory
	}
	output, err := cmd.Output()
	return output, err
}
