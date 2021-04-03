package builder

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/io"
)

type Orchestrator interface {
	RunBuild(specifiedLambdas []string) error
}

type orchestrator struct {
	builder    Builder
	filesystem io.Filesystem
	out        io.Printer
	lambdasDir string
}

func NewOrchestrator(builder Builder, filesystem io.Filesystem, out io.Printer) Orchestrator {
	return &orchestrator{
		builder:    builder,
		filesystem: filesystem,
		out: out,
		lambdasDir: "lambdas",
	}
}

func (o *orchestrator) RunBuild(specifiedLambdas []string) error {
	allLambdas, err := o.findLambdaNames()
	if err != nil {
		return fmt.Errorf("Unable to find directory %s", o.lambdasDir)
	}

	var lambdasToBuild []string
	if len(specifiedLambdas) > 0 {
		for _, lambda := range specifiedLambdas {
			if !contains(lambda, allLambdas) {
				return fmt.Errorf("Unable to find lambda %s", lambda)
			}
		}
		lambdasToBuild = specifiedLambdas
	} else {
		lambdasToBuild = allLambdas
	}

	for _, lambda := range lambdasToBuild {
		o.out.Printlnf("🔨 Building %s...", lambda)
		err := o.builder.BuildLambda(lambda)
		if err != nil {
			return fmt.Errorf("Error building %s", lambda)
		}
	}
	o.out.Printlnf("✅ Done")
	return nil
}

func (o *orchestrator) findLambdaNames() (lambdaNames []string, err error) {
	dirContents, err := o.filesystem.ReadDir(o.lambdasDir)
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
