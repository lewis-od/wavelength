package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/lambda-build/pkg/builder"
	"os"
)

type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Uploader struct {
	client S3PutObjectAPI
	ctx    context.Context
}

func NewS3Uploader(client S3PutObjectAPI, ctx context.Context) builder.Uploader {
	return &S3Uploader{
		client: client,
		ctx:    ctx,
	}
}

func (s *S3Uploader) UploadLambda(version, bucketName, lambdaName, artifactLocation string) error {
	uploadLocation := fmt.Sprintf("%s/%s.zip", version, lambdaName)

	file, err := os.Open(artifactLocation)
	if err != nil {
		return err
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &uploadLocation,
		Body:   file,
	}
	_, err = s.client.PutObject(s.ctx, input)
	return err
}
