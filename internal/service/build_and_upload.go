package service

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/builder"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/lewis-od/lambda-build/internal/terraform"
)

type BuildAndUploadService interface {
	Run(version string, lambdas []string, skipBuild bool)
}

type buildAndUploadService struct {
	orchestrator       builder.Orchestrator
	terraform          terraform.Terraform
	filesystem         io.Filesystem
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
	filesystem io.Filesystem,
	out io.Printer,
) BuildAndUploadService {
	return &buildAndUploadService{
		orchestrator:       orchestrator,
		terraform:          terraform,
		filesystem:         filesystem,
		out:                out,
		artifactDeployment: "terraform/deployments/artifact-storage",
		lambdasDir:         "lambdas",
	}
}

func (c *buildAndUploadService) Run(version string, lambdas []string, skipBuild bool) {
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

func (c *buildAndUploadService) validateLambdaNames(providedNames []string) ([]string, error) {
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

func (c *buildAndUploadService) findLambdaNames() (lambdaNames []string, err error) {
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
