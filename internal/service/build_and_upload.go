package service

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/builder"
	"github.com/lewis-od/lambda-build/internal/finder"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/lewis-od/lambda-build/internal/terraform"
)

type BuildAndUploadService interface {
	Run(version string, lambdas []string, skipBuild bool)
}

type buildAndUploadService struct {
	orchestrator       builder.Orchestrator
	terraform          terraform.Terraform
	finder             finder.Finder
	out                io.Printer
	artifactDeployment string
	lambdasDir         string
}

type uploadArguments struct {
	version string
	lambdas []string
}

func NewBuildAndUploadService(
	orchestrator builder.Orchestrator,
	terraform terraform.Terraform,
	finder finder.Finder,
	out io.Printer,
) BuildAndUploadService {
	return &buildAndUploadService{
		orchestrator:       orchestrator,
		terraform:          terraform,
		finder:             finder,
		out:                out,
		artifactDeployment: "terraform/deployments/artifact-storage",
		lambdasDir:         "lambdas",
	}
}

func (c *buildAndUploadService) Run(version string, lambdas []string, skipBuild bool) {
	lambdasToUpload, err := c.finder.GetLambdas(lambdas)
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

func (c *buildAndUploadService) findArtifactBucketName() (string, error) {
	outputs, err := c.terraform.Output(c.artifactDeployment)
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
