package builder

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/io"
	"github.com/lewis-od/wavelength/internal/progress"
)

type Orchestrator interface {
	BuildLambdas(lambdas []string) []*BuildResult
	UploadLambdas(version, bucketName string, lambdas []string) []*BuildResult
}

type orchestrator struct {
	builder  Builder
	uploader Uploader
	display  progress.BuildDisplay
	out      io.Printer
}

func NewOrchestrator(
	builder Builder,
	uploader Uploader,
	display progress.BuildDisplay,
	out io.Printer) Orchestrator {
	return &orchestrator{
		builder:  builder,
		uploader: uploader,
		display:  display,
		out:      out,
	}
}

func (o *orchestrator) BuildLambdas(lambdas []string) []*BuildResult {
	resultChan := make(chan *BuildResult)

	for _, lambda := range lambdas {
		o.display.Started(lambda)
	}
	for _, lambda := range lambdas {
		go func(lambdaName string) {
			resultChan <- o.builder.BuildLambda(lambdaName)
		}(lambda)
	}

	results := make([]*BuildResult, 0, len(lambdas))
	for {
		result := <-resultChan
		results = append(results, result)
		o.display.Completed(result.LambdaName, result.Error == nil)
		if len(results) == len(lambdas) {
			break
		}
	}

	return getFailedResults(results)
}

func (o *orchestrator) UploadLambdas(version, bucketName string, lambdas []string) []*BuildResult {
	resultChan := make(chan *BuildResult)

	for _, lambda := range lambdas {
		o.out.Printlnf("☁️  Uploading %s...", lambda)
		go func(lambdaName string) {
			artifact := fmt.Sprintf("lambdas/%s/dist/%s.zip", lambdaName, lambdaName)
			resultChan <- o.uploader.UploadLambda(version, bucketName, lambdaName, artifact)
		}(lambda)
	}

	results := make([]*BuildResult, 0, len(lambdas))
	for {
		result := <-resultChan
		results = append(results, result)
		if len(results) == len(lambdas) {
			break
		}
	}

	return getFailedResults(results)
}

func getFailedResults(results []*BuildResult) []*BuildResult {
	failedResults := make([]*BuildResult, 0, len(results))
	for _, result := range results {
		if result.Error != nil {
			failedResults = append(failedResults, result)
		}
	}
	return failedResults
}
