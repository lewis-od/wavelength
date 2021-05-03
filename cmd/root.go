package cmd

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/find"
	"github.com/lewis-od/wavelength/internal/ports/aws"
	"github.com/lewis-od/wavelength/internal/ports/lerna"
	"github.com/lewis-od/wavelength/internal/ports/stdout"
	"github.com/lewis-od/wavelength/internal/ports/system"
	"github.com/lewis-od/wavelength/internal/terraform"
	"github.com/spf13/cobra"
)

// To be loaded from config
const projectName = "jarvis"
const lambdasDir = "lambdas"
const artifactStorageComponent = "terraform/deployments/artifact-storage"

var printer = stdout.NewPrinter()
var lernaBuilder = lerna.NewLerna(system.NewExecutor("lerna"), "jarvis")
var awsContext = context.Background()
var lambdaUploader = aws.NewS3Uploader(newS3Client(awsContext), awsContext)
var orchestrator = builder.NewOrchestrator(lernaBuilder, lambdaUploader, printer)
var tfExec = terraform.NewTerraform(system.NewExecutor("terraform"))
var filesystem = system.NewFilesystem()
var finder = find.NewLambdaFinder(filesystem, tfExec, lambdasDir, artifactStorageComponent)
var updater = aws.NewLambdaUpdater(newLambdaClient(awsContext), awsContext)

func newS3Client(ctx context.Context) *s3.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return s3.NewFromConfig(cfg)
}

func newLambdaClient(ctx context.Context) *lambda.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	return lambda.NewFromConfig(cfg)
}

var rootCmd = &cobra.Command{
	Use:   "wavelength",
	Short: "Opinionated tool for building and deploying lambdas using Terraform & Node.js",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()
}
