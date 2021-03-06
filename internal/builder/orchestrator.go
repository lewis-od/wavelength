package builder

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/io"
	"github.com/lewis-od/wavelength/internal/progress"
)

type Orchestrator interface {
	BuildLambdas(lambdas []string) []*BuildResult
	UploadLambdas(version, bucketName string, lambdas []string, role *Role) []*BuildResult
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

	o.display.Init(progress.Build)
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

func (o *orchestrator) UploadLambdas(version, bucketName string, lambdas []string, role *Role) []*BuildResult {
	resultChan := make(chan *BuildResult)

	o.display.Init(progress.Upload)
	for _, lambda := range lambdas {
		o.display.Started(lambda)
	}
	for _, lambda := range lambdas {
		go func(lambdaName string) {
			artifact := fmt.Sprintf("lambdas/%s/dist/%s.zip", lambdaName, lambdaName)
			resultChan <- o.uploader.UploadLambda(version, bucketName, lambdaName, artifact, role)
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

func getFailedResults(results []*BuildResult) []*BuildResult {
	failedResults := make([]*BuildResult, 0, len(results))
	for _, result := range results {
		if result.Error != nil {
			failedResults = append(failedResults, result)
		}
	}
	return failedResults
}
