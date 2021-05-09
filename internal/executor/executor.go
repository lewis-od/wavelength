package executor

type CommandContext struct {
	Directory string
}

type CommandExecutor interface {
	Execute(args []string) error
	ExecuteAndCapture(args []string, context *CommandContext) ([]byte, error)
}
