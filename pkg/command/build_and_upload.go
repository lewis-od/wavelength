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
	out               io.Printer
	artifactWorkspace string
}

func NewBuildAndUploadCommand(
	orchestrator builder.Orchestrator,
	terraform terraform.Terraform,
	out io.Printer,
) *BuildAndUploadCommand {
	return &BuildAndUploadCommand{
		orchestrator:      orchestrator,
		terraform:         terraform,
		out:               out,
		artifactWorkspace: "terraform/deployments/artifact-storage",
	}
}

func (c *BuildAndUploadCommand) Name() string {
	return "upload"
}

func (c *BuildAndUploadCommand) Description() string {
	return "Build and upload to S3"
}

func (c *BuildAndUploadCommand) Run(arguments []string) {
	err := c.orchestrator.RunBuild(arguments)
	if err != nil {
		c.out.Println("❌", err)
		return
	}

	bucketName, err := c.findArtifactBucketName()
	if err != nil {
		c.out.Println(err)
		return
	}
	c.out.Println("Uploading to", bucketName)
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
