package builder

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/io"
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
	for _, lambda := range lambdas {
		o.out.Printlnf("🔨 Building %s...", lambda)
		err := o.builder.BuildLambda(lambda)
		if err != nil {
			return fmt.Errorf("Error building %s", lambda)
		}
	}
	o.out.Printlnf("✅ Build complete")
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
