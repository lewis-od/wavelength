package find

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/lewis-od/lambda-build/internal/terraform"
)

type Finder interface {
	FindLambdas(lambdaNames []string) ([]string, error)
	FindArtifactBucketName() (string, error)
}

type lambdaFinder struct {
	filesystem                           io.Filesystem
	tf                                   terraform.Terraform
	lambdasDir, artifactStorageComponent string
}

func NewLambdaFinder(filesystem io.Filesystem,
	tf terraform.Terraform,
	lambdasDir, artifactStorageComponent string) Finder {
	return &lambdaFinder{
		filesystem:               filesystem,
		tf:                       tf,
		lambdasDir:               lambdasDir,
		artifactStorageComponent: artifactStorageComponent,
	}
}

func (f *lambdaFinder) FindLambdas(providedNames []string) ([]string, error) {
	allLambdas, err := f.findLambdaNames()
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

func (f *lambdaFinder) findLambdaNames() (lambdaNames []string, err error) {
	dirContents, err := f.filesystem.ReadDir(f.lambdasDir)
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

func (f *lambdaFinder) FindArtifactBucketName() (string, error) {
	outputs, err := f.tf.Output(f.artifactStorageComponent)
	if err != nil {
		return "", err
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
