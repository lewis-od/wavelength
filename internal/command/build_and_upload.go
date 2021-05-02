package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/builder"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/lewis-od/lambda-build/internal/terraform"
)

type BuildAndUploadCommand interface {
	Run(version string, lambdas []string, skipBuild bool)
}

type buildAndUploadCommand struct {
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
) BuildAndUploadCommand {
	return &buildAndUploadCommand{
		orchestrator:      orchestrator,
		terraform:         terraform,
		filesystem:        filesystem,
		out:               out,
		artifactWorkspace: "terraform/deployments/artifact-storage",
		lambdasDir:        "lambdas",
	}
}

func (c *buildAndUploadCommand) Run(version string, lambdas []string, skipBuild bool) {
	lambdasToUpload, err := c.validateLambdaNames(lambdas)
	if err != nil {
		c.out.PrintErr(err)
		return
	}
	c.out.Printlnf("🏗  Orchestrating upload of version %s of %s", version, lambdasToUpload)

	bucketName, err := c.findArtifactBucketName()
	if err != nil {
		c.out.PrintErr(err)
		return
	}
	c.out.Printlnf("🪣 Found artifact bucket %s", bucketName)

	if !skipBuild {
		err = c.orchestrator.BuildLambdas(lambdasToUpload)
		if err != nil {
			c.out.PrintErr(err)
			return
		}
	}
	err = c.orchestrator.UploadLambdas(version, bucketName, lambdasToUpload)
	if err != nil {
		c.out.PrintErr(err)
		return
	}
}

func (c *buildAndUploadCommand) findArtifactBucketName() (string, error) {
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

func (c *buildAndUploadCommand) validateLambdaNames(providedNames []string) ([]string, error) {
	allLambdas, err := c.findLambdaNames()
	if err != nil {
		return nil, err
	}

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

	return toBuild, nil
}

func (c *buildAndUploadCommand) findLambdaNames() (lambdaNames []string, err error) {
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
