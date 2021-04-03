package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/lerna"
	"github.com/lewis-od/lambda-build/pkg/terraform"
)

type BuildAndUploadCommand struct {
	lerna             lerna.Lerna
	terraform         terraform.Terraform
	filesystem        Filesystem
	lambdasDirectory  string
	artifactWorkspace string
}

func NewBuildAndUploadCommand(
	lerna lerna.Lerna,
	terraform terraform.Terraform,
	filesystem Filesystem,
) *BuildAndUploadCommand {
	return &BuildAndUploadCommand{
		lerna:             lerna,
		terraform:         terraform,
		filesystem:        filesystem,
		lambdasDirectory:  "lambdas",
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
	lambdaNames, err := c.findLambdaNames()
	if err != nil {
		fmt.Printf("Unable to read directory %s\n", c.lambdasDirectory)
		return
	}

	if len(arguments) == 0 {
		err = c.buildLambdas(lambdaNames)
	} else {
		lambdaName := arguments[0]
		if !contains(lambdaName, lambdaNames) {
			fmt.Printf("Lambda not found: %s\n", lambdaName)
			return
		}
		err = c.buildLambdas([]string{lambdaName})
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	bucketName, err := c.findArtifactBucketName()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Uploading to ", bucketName)
}

func (c *BuildAndUploadCommand) findLambdaNames() (lambdaNames []string, err error) {
	dirContents, err := c.filesystem.ReadDir(c.lambdasDirectory)
	if err != nil {
		return
	}
	for _, lambdaDir := range dirContents {
		if lambdaDir.IsDir() {
			lambdaNames = append(lambdaNames, lambdaDir.Name())
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

func (c *BuildAndUploadCommand) buildLambdas(names []string) error {
	for _, lambdaName := range names {
		fmt.Printf("🔨 Building %s lambda...\n", lambdaName)
		err := c.lerna.BuildLambda(lambdaName)
		if err != nil {
			return fmt.Errorf("❌ Error building %s\n%s\n", lambdaName, err)
		}
		artifactPath := fmt.Sprintf("%s/%s/dist/%s.zip", c.lambdasDirectory, lambdaName, lambdaName)
		if !c.filesystem.FileExists(artifactPath) {
			return fmt.Errorf("❌ Artifact %s not found, did the build succeed?", artifactPath)
		}
	}
	fmt.Println("✅ Done")
	return nil
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
