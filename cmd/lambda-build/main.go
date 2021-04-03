package main

import (
	"github.com/lewis-od/lambda-build/pkg/io/system"
	"github.com/lewis-od/lambda-build/pkg/lerna"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"os"

	"github.com/lewis-od/lambda-build/pkg/command"
)

func main() {
	config := command.CLIConfig{
		Name:    "lambda-build",
		Version: "0.1",
	}
	app := command.NewCLI(config)

	lernaExec := lerna.NewLerna(system.NewExecutor("lerna"), "jarvis")
	tfExec := terraform.NewTerraform(system.NewExecutor("terraform"))
	filesystem := &system.OSFilesystem{}
	buildCommand := command.NewBuildAndUploadCommand(lernaExec, tfExec, filesystem)
	app.AddCommand(buildCommand)

	app.Run(os.Args)
}
