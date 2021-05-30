package service

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/find"
	"github.com/lewis-od/wavelength/internal/io"
)

type BuildAndUploadService interface {
	Run(version string, lambdas []string, skipBuild bool, role *builder.Role)
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

func (c *buildAndUploadService) Run(version string, lambdas []string, skipBuild bool, role *builder.Role) {
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

	uploadErrors := c.orchestrator.UploadLambdas(version, bucketName, lambdasToUpload, role)
	if len(uploadErrors) != 0 {
		c.printUploadErrors(uploadErrors)
		return
	}

	c.out.Println("✅ Done!")
}

func (c *buildAndUploadService) printBuildErrors(buildResults []*builder.BuildResult) {
	c.printErrors(build, buildResults)
}

func (c *buildAndUploadService) printUploadErrors(buildResults []*builder.BuildResult) {
	c.printErrors(upload, buildResults)
}

const build = "building"
const upload = "uploading"

func (c *buildAndUploadService) printErrors(action string, buildResults []*builder.BuildResult) {
	for _, result := range buildResults {
		err := fmt.Errorf("Error %s lambda %s\n%s\n", action, result.LambdaName, result.Output)
		c.out.PrintErr(err)
	}
}
