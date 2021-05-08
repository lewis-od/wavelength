package builder

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/error_logger"
	"github.com/lewis-od/wavelength/internal/io"
)

type Orchestrator interface {
	BuildLambdas(lambdas []string) error
	UploadLambdas(version, bucketName string, lambdas []string) error
}

type orchestrator struct {
	builder   Builder
	uploader  Uploader
	errLogger error_logger.ErrorLogger
	out       io.Printer
}

func NewOrchestrator(
	builder Builder,
	uploader Uploader,
	errLogger error_logger.ErrorLogger,
	out io.Printer) Orchestrator {
	return &orchestrator{
		builder:   builder,
		uploader:  uploader,
		errLogger: errLogger,
		out:       out,
	}
}

func (o *orchestrator) BuildLambdas(lambdas []string) error {
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

	failedLambdas := make([]string, 0, len(lambdas))
	for _, result := range results {
		if result.Error != nil {
			failedLambdas = append(failedLambdas, result.LambdaName)
			o.errLogger.AddError(mapToError(result))
		}
	}
	if len(failedLambdas) != 0 {
		err := o.errLogger.WriteLogFile()
		if err != nil {
			return err
		}
		return fmt.Errorf("Error building lambdas %s", failedLambdas)
	}

	o.out.Println("✅ Build complete")
	return nil
}

func mapToError(result *BuildResult) *error_logger.WavelengthError {
	return &error_logger.WavelengthError{
		Lambda: result.LambdaName,
		Output: result.Output,
	}
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
