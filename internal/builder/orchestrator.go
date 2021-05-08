package builder

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/io"
)

type Orchestrator interface {
	BuildLambdas(lambdas []string) error
	UploadLambdas(version, bucketName string, lambdas []string) error
}

type orchestrator struct {
	builder  Builder
	uploader Uploader
	out      io.Printer
}

func NewOrchestrator(builder Builder, uploader Uploader, out io.Printer) Orchestrator {
	return &orchestrator{
		builder:  builder,
		uploader: uploader,
		out:      out,
	}
}

func (o *orchestrator) BuildLambdas(lambdas []string) error {
	errChan := make(chan error)
	successChan := make(chan bool)

	for _, lambda := range lambdas {
		o.out.Printlnf("🔨 Building %s...", lambda)
		go func(lambdaName string) {
			_, err := o.builder.BuildLambda(lambdaName)
			if err != nil {
				errChan <- fmt.Errorf("Error building %s", lambdaName)
			} else {
				successChan <- true
			}
		}(lambda)
	}

	errs := make([]error, 0, len(lambdas))
	completedBuilds := 0
	for {
		select {
		case err := <-errChan:
			errs = append(errs, err)
			completedBuilds++
		case <-successChan:
			completedBuilds++
		}
		if completedBuilds == len(lambdas) {
			break
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("%s", errs)
	}

	o.out.Println("✅ Build complete")
	return nil
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
