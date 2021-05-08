package service

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/find"
	"github.com/lewis-od/wavelength/internal/io"
)

type BuildAndUploadService interface {
	Run(version string, lambdas []string, skipBuild bool)
}

type buildAndUploadService struct {
	orchestrator builder.Orchestrator
	finder       find.Finder
	out          io.Printer
}

func NewBuildAndUploadService(
	orchestrator builder.Orchestrator,
	finder find.Finder,
	out io.Printer,
) BuildAndUploadService {
	return &buildAndUploadService{
		orchestrator: orchestrator,
		finder:       finder,
		out:          out,
	}
}

func (c *buildAndUploadService) Run(version string, lambdas []string, skipBuild bool) {
	lambdasToUpload, err := c.finder.FindLambdas(lambdas)
	if err != nil {
		c.out.PrintErr(err)
		return
	}
	c.out.Printlnf("🏗  Orchestrating upload of version %s of %s", version, lambdasToUpload)

	bucketName, err := c.finder.FindArtifactBucketName()
	if err != nil {
		c.out.PrintErr(err)
		return
	}
	c.out.Printlnf("🪣 Found artifact bucket %s", bucketName)

	if !skipBuild {
		failedBuilds := c.orchestrator.BuildLambdas(lambdasToUpload)
		if len(failedBuilds) != 0 {
			c.printBuildErrors(failedBuilds)
			return
		}
	}
	err = c.orchestrator.UploadLambdas(version, bucketName, lambdasToUpload)
	if err != nil {
		c.out.PrintErr(err)
		return
	}
}

func (c *buildAndUploadService) printBuildErrors(buildResults []*builder.BuildResult) {
	for _, result := range buildResults {
		err := fmt.Errorf("Error building lambda %s\n%s\n%s\n", result.LambdaName, result.Error, result.Output)
		c.out.PrintErr(err)
	}
}
