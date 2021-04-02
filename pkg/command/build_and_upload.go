package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/lerna"
)

type BuildAndUploadCommand struct {
	lerna lerna.Lerna
	filesystem Filesystem
}

func NewBuildAndUploadCommand(lerna lerna.Lerna, filesystem Filesystem) *BuildAndUploadCommand {
	return &BuildAndUploadCommand{
		lerna: lerna,
		filesystem: filesystem,
	}
}

func (c *BuildAndUploadCommand) Name() string {
	return "upload"
}

func (c *BuildAndUploadCommand) Description() string {
	return "Build and upload to S3"
}

func (c *BuildAndUploadCommand) Run(arguments []string) {
	lambdasDir := "lambdas"
	lambdaNames, err := c.findLambdaNames(lambdasDir)
	if err != nil {
		fmt.Printf("Unable to read directory %s\n", lambdasDir)
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
	}
	fmt.Println("✅ Done")
	return nil
}

func (c *BuildAndUploadCommand) findLambdaNames(lambdasDir string) (lambdaNames []string, err error) {
	dirContents, err := c.filesystem.ReadDir(lambdasDir)
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
