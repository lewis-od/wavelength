package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/lambda-build/pkg/builder"
	"github.com/lewis-od/lambda-build/pkg/ports/aws"
	"github.com/lewis-od/lambda-build/pkg/ports/lerna"
	"github.com/lewis-od/lambda-build/pkg/ports/stdout"
	"github.com/lewis-od/lambda-build/pkg/ports/system"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"github.com/spf13/cobra"
)

var printer = stdout.NewPrinter()
var lernaBuilder = lerna.NewLerna(system.NewExecutor("lerna"), "jarvis")
var awsContext = context.Background()
var lambdaUploader = aws.NewS3Uploader(newS3Client(awsContext), awsContext)
var orchestrator = builder.NewOrchestrator(lernaBuilder, lambdaUploader, printer)
var tfExec = terraform.NewTerraform(system.NewExecutor("terraform"))
var filesystem = system.NewFilesystem()

func newS3Client(ctx context.Context) *s3.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg)
}

var rootCmd = &cobra.Command{
	Use:   "lambda-build",
	Short: "Tool for working with lambdas",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()
}
