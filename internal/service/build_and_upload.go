package service

import (
	"github.com/lewis-od/lambda-build/internal/builder"
	"github.com/lewis-od/lambda-build/internal/find"
	"github.com/lewis-od/lambda-build/internal/io"
)

type BuildAndUploadService interface {
	Run(version string, lambdas []string, skipBuild bool)
}

type buildAndUploadService struct {
	orchestrator       builder.Orchestrator
	finder             find.Finder
	out                io.Printer
}

type uploadArguments struct {
	version string
	lambdas []string
}

func NewBuildAndUploadService(
	orchestrator builder.Orchestrator,
	finder find.Finder,
	out io.Printer,
) BuildAndUploadService {
	return &buildAndUploadService{
		orchestrator:       orchestrator,
		finder:             finder,
		out:                out,
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
