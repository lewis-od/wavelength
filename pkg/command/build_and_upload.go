package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/builder"
	"github.com/lewis-od/lambda-build/pkg/io"
	"github.com/lewis-od/lambda-build/pkg/terraform"
)

type BuildAndUploadCommand struct {
	orchestrator      builder.Orchestrator
	terraform         terraform.Terraform
	filesystem        io.Filesystem
	out               io.Printer
	artifactWorkspace string
	lambdasDir        string
}

type uploadArguments struct {
	version string
	lambdas []string
}

func NewBuildAndUploadCommand(
	orchestrator builder.Orchestrator,
	terraform terraform.Terraform,
	filesystem io.Filesystem,
	out io.Printer,
) *BuildAndUploadCommand {
	return &BuildAndUploadCommand{
		orchestrator:      orchestrator,
		terraform:         terraform,
		filesystem:        filesystem,
		out:               out,
		artifactWorkspace: "terraform/deployments/artifact-storage",
		lambdasDir:        "lambdas",
	}
}

func (c *BuildAndUploadCommand) Name() string {
	return "upload"
}

func (c *BuildAndUploadCommand) Description() string {
	return "Build and upload to S3"
}

func (c *BuildAndUploadCommand) Run(args []string) {
	arguments, err := c.parseArguments(args)
	if err != nil {
		c.printError(err)
		return
	}
	c.out.Printlnf("🏗  Building version %s of %s", arguments.version, arguments.lambdas)

	bucketName, err := c.findArtifactBucketName()
	if err != nil {
		c.printError(err)
		return
	}
	c.out.Printlnf("🪣 Found artifact bucket %s", bucketName)

	err = c.orchestrator.RunBuild(arguments.version, bucketName, arguments.lambdas)
	if err != nil {
		c.printError(err)
		return
	}
}

func (c *BuildAndUploadCommand) printError(err error) {
	c.out.Println("❌", err)
}

func (c *BuildAndUploadCommand) findArtifactBucketName() (string, error) {
	outputs, err := c.terraform.Output(c.artifactWorkspace)
	if err != nil {
		return "", fmt.Errorf("Could not determine name of artifact bucket from tf state\n%s", err)
	}
	bucketName, outputExists := outputs["bucket_name"]
	if !outputExists {
		outputNames := make([]string, 0, len(outputs))
		for output := range outputs {
			outputNames = append(outputNames, output)
		}
		return "", fmt.Errorf("No output named bucket_name found in %s", outputNames)
	}
	return bucketName.Value, nil
}

func (c *BuildAndUploadCommand) parseArguments(args []string) (*uploadArguments, error) {
	allLambdas, err := c.findLambdaNames()
	if err != nil {
		return nil, err
	}

	if len(args) < 1 {
		return &uploadArguments{}, fmt.Errorf("Version not specified")
	}
	version := args[0]

	providedNames := args[1:]
	toBuild := make([]string, 0, len(allLambdas))
	if len(providedNames) == 0 {
		toBuild = allLambdas
	} else {
		for _, name := range providedNames {
			if !contains(name, allLambdas) {
				return nil, fmt.Errorf("Could not find lambda %s", name)
			}
			toBuild = append(toBuild, name)
		}
	}

	return &uploadArguments{
		version: version,
		lambdas: toBuild,
	}, nil
}

func (c *BuildAndUploadCommand) findLambdaNames() (lambdaNames []string, err error) {
	dirContents, err := c.filesystem.ReadDir(c.lambdasDir)
	if err != nil {
		return
	}
	for _, lambdaDir := range dirContents {
		if lambdaDir.IsDir {
			lambdaNames = append(lambdaNames, lambdaDir.Name)
		}
	}
	return
}

func contains(target string, items []string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
