package builder

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/io"
)

type Orchestrator interface {
	BuildLambdas(lambdas []string) []*BuildResult
	UploadLambdas(version, bucketName string, lambdas []string) error
}

type orchestrator struct {
	builder  Builder
	uploader Uploader
	out      io.Printer
}

func NewOrchestrator(
	builder Builder,
	uploader Uploader,
	out io.Printer) Orchestrator {
	return &orchestrator{
		builder:  builder,
		uploader: uploader,
		out:      out,
	}
}

func (o *orchestrator) BuildLambdas(lambdas []string) []*BuildResult {
	resultChan := make(chan *BuildResult)

	for _, lambda := range lambdas {
		o.out.Printlnf("🔨 Building %s...", lambda)
		go func(lambdaName string) {
			resultChan <- o.builder.BuildLambda(lambdaName)
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

	failedLambdas := make([]*BuildResult, 0, len(lambdas))
	for _, result := range results {
		if result.Error != nil {
			failedLambdas = append(failedLambdas, result)
		}
	}
	if len(failedLambdas) == 0 {
		o.out.Println("✅ Build complete")
	}
	return failedLambdas
}

func (o *orchestrator) UploadLambdas(version, bucketName string, lambdas []string) error {
	for _, lambda := range lambdas {
		artifact := fmt.Sprintf("lambdas/%s/dist/%s.zip", lambda, lambda)
		o.out.Printlnf("☁️  Uploading %s...", lambda)
		err := o.uploader.UploadLambda(version, bucketName, lambda, artifact)
		if err != nil {
			return fmt.Errorf("Error uploading lambda %s\n%s", lambda, err)
		}
	}
	o.out.Println("✅ Upload complete")
	return nil
}
