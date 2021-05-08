package cmd

import (
	"context"
	"fmt"
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
	"github.com/spf13/viper"
	"os"
)

var configFile string

// Loaded from config
var projectName string
var artifactStorageComponent string
var bucketOutputName string
var lambdasDir string

var printer = stdout.NewPrinter()
var lernaBuilder = lerna.NewLerna(system.NewExecutor("lerna"), &projectName)
var awsContext = context.Background()
var lambdaUploader = aws.NewS3Uploader(newS3Client(awsContext), awsContext)
var tfExec = terraform.NewTerraform(system.NewExecutor("terraform"))
var filesystem = system.NewFilesystem()
var orchestrator = builder.NewOrchestrator(lernaBuilder, lambdaUploader, printer)
var finder = find.NewLambdaFinder(filesystem, tfExec, &lambdasDir, &artifactStorageComponent, &bucketOutputName)
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file to use")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			printer.PrintErr(err)
			os.Exit(1)
		}

		viper.AddConfigPath(cwd)
		viper.SetConfigName(".wavelength")
	}

	viper.SetDefault("lambdas", "lambdas")

	if err := viper.ReadInConfig(); err != nil {
		printer.PrintErr(err)
		os.Exit(1)
	}

	setFromConfig(&projectName, "projectName", true)
	setFromConfig(&artifactStorageComponent, "artifactStorage.terraformDir", true)
	setFromConfig(&bucketOutputName, "artifactStorage.outputName", true)
	setFromConfig(&lambdasDir, "lambdas", false)
}

func setFromConfig(holder *string, key string, required bool) {
	value := viper.GetString(key)
	if value == "" && required {
		err := fmt.Errorf("value %s not found in config", key)
		printer.PrintErr(err)
		os.Exit(1)
	}
	*holder = value
}
