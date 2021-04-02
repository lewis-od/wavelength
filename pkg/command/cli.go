package command

import (
	"fmt"
)

type Command interface {
	Name() string
	Description() string
	Run(arguments []string)
}

type CLIConfig struct {
	Name    string
	Version string
}

type appCLI struct {
	Name     string
	Version  string
	commands map[string]Command
}

func NewCLI(config CLIConfig) *appCLI {
	return &appCLI{
		Name:     config.Name,
		Version:  config.Version,
		commands: make(map[string]Command),
	}
}

func (cli *appCLI) AddCommand(command Command) {
	cli.commands[command.Name()] = command
}

func (cli *appCLI) Run(args []string) {
	cliPath := args[0]
	if len(args) == 1 {
		cli.printUsage(cliPath)
		return
	}

	commandName := args[1]
	command := cli.commands[commandName]
	if command != nil {
		command.Run(args[2:])
	} else {
		fmt.Printf("Command not found %s\n", commandName)
		cli.printUsage(cliPath)
	}
}

func (cli *appCLI) printUsage(commandPath string) {
	fmt.Printf("Usage: %s [command] [arguments]\n", commandPath)
	fmt.Println("Where [command] is one of:")
	for name, command := range cli.commands {
		fmt.Printf("\t%s - %s\n", name, command.Description())
	}
}
