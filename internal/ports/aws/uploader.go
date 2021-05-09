package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/wavelength/internal/builder"
	"os"
)

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type s3Uploader struct {
	client S3PutObjectAPI
	ctx    context.Context
}

func NewS3Uploader(client S3PutObjectAPI, ctx context.Context) builder.Uploader {
	return &s3Uploader{
		client: client,
		ctx:    ctx,
	}
}

func (s *s3Uploader) UploadLambda(version, bucketName, lambdaName, artifactLocation string) *builder.BuildResult {
	uploadLocation := fmt.Sprintf("%s/%s.zip", version, lambdaName)

	file, err := os.Open(artifactLocation)
	if err != nil {
		return &builder.BuildResult{
			LambdaName: lambdaName,
			Error:      err,
			Output:     []byte(err.Error()),
		}
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &uploadLocation,
		Body:   file,
	}
	_, err = s.client.PutObject(s.ctx, input)
	output := ""
	if err != nil {
		output = err.Error()
	}
	return &builder.BuildResult{
		LambdaName: lambdaName,
		Error: err,
		Output: []byte(output),
	}
}
