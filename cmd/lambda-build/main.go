package main

import (
	"fmt"
	"os"

	"github.com/lewis-od/lambda-build/pkg/command"
)

type HelloCommand struct{}

func (h *HelloCommand) Name() string {
	return "hello"
}

func (h *HelloCommand) Run(arguments []string) {
	fmt.Printf("Hello with arguments: %s\n", arguments)
}

func (h *HelloCommand) Description() string {
	return "Say hello"
}

func main() {
	config := command.CLIConfig{
		Name:    "lambda-build",
		Version: "0.1",
	}
	app := command.NewCLI(config)
	app.AddCommand(&HelloCommand{})
	app.Run(os.Args)
}
