package main

import (
	"github.com/lewis-od/lambda-build/pkg/builder"
	"github.com/lewis-od/lambda-build/pkg/command"
	"github.com/lewis-od/lambda-build/pkg/ports/lerna"
	"github.com/lewis-od/lambda-build/pkg/ports/stdout"
	"github.com/lewis-od/lambda-build/pkg/ports/system"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"os"
)

func main() {
	config := command.CLIConfig{
		Name:    "lambda-build",
		Version: "0.1",
	}
	app := command.NewCLI(config)

	printer := stdout.NewPrinter()
	filesystem := &system.OSFilesystem{}
	lernaBuilder := lerna.NewLerna(system.NewExecutor("lerna"), "jarvis")
	orchestrator := builder.NewOrchestrator(lernaBuilder, filesystem, printer)
	tfExec := terraform.NewTerraform(system.NewExecutor("terraform"))
	buildCommand := command.NewBuildAndUploadCommand(orchestrator, tfExec, printer)

	app.AddCommand(buildCommand)

	app.Run(os.Args)
}
