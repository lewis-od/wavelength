package finder

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/io"
)

type Finder interface {
	GetLambdas(lambdaNames []string) ([]string, error)
}

type lambdaFinder struct {
	filesystem io.Filesystem
	lambdasDir string
}

func NewLambdaFinder(filesystem io.Filesystem, lambdasDir string) Finder {
	return &lambdaFinder{
		filesystem: filesystem,
		lambdasDir: lambdasDir,
	}
}

func (f *lambdaFinder) GetLambdas(providedNames []string) ([]string, error) {
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
