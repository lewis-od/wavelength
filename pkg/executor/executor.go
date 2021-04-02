package executor

type CommandContext struct {
	Directory string
}

type CommandExecutor interface {
	ExecuteWithContext(args []string, context *CommandContext) ([]byte, error)
}
