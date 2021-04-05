package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/lambda-build/pkg/builder"
	"github.com/lewis-od/lambda-build/pkg/command"
	"github.com/lewis-od/lambda-build/pkg/ports/aws"
	"github.com/lewis-od/lambda-build/pkg/ports/lerna"
	"github.com/lewis-od/lambda-build/pkg/ports/stdout"
	"github.com/lewis-od/lambda-build/pkg/ports/system"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"os"
)

func main() {
	cliConfig := command.CLIConfig{
		Name:    "lambda-build",
		Version: "0.1",
	}
	app := command.NewCLI(cliConfig)

	printer := stdout.NewPrinter()
	lernaBuilder := lerna.NewLerna(system.NewExecutor("lerna"), "jarvis")
	awsContext := context.Background()
	lambdaUploader := aws.NewS3Uploader(newS3Client(awsContext), awsContext)
	orchestrator := builder.NewOrchestrator(lernaBuilder, lambdaUploader, printer)
	tfExec := terraform.NewTerraform(system.NewExecutor("terraform"))
	filesystem := system.NewFilesystem()
	buildCommand := command.NewBuildAndUploadCommand(orchestrator, tfExec, filesystem, printer)

	app.AddCommand(buildCommand)

	app.Run(os.Args)
}

func newS3Client(ctx context.Context) *s3.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg)
}
