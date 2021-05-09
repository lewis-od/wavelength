package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockS3Client struct {
	mock.Mock
	input *s3.PutObjectInput
}

func (m *mockS3Client) PutObject(ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.input = params
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func TestS3Uploader_UploadLambda_Success(t *testing.T) {
	mockClient := new(mockS3Client)
	mockClient.On(
		"PutObject",
		context.TODO(), mock.Anything, mock.Anything, mock.Anything,
	).Return(&s3.PutObjectOutput{}, nil)

	uploader := NewS3Uploader(mockClient, context.TODO())

	result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip")

	assert.Nil(t, result.Error)
	mockClient.AssertExpectations(t)
	assert.Equal(t, "test/my-lambda.zip", *mockClient.input.Key)
}

func TestS3Uploader_UploadLambda_FileOpenError(t *testing.T) {
	mockClient := new(mockS3Client)

	uploader := NewS3Uploader(mockClient, context.TODO())

	result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/missing.zip")

	assert.NotNil(t, result.Error)
	mockClient.AssertExpectations(t)
}

func TestS3Uploader_UploadLambda_S3PutError(t *testing.T) {
	mockClient := new(mockS3Client)
	uploadError := fmt.Errorf("upload error")
	mockClient.On(
		"PutObject",
		context.TODO(), mock.Anything, mock.Anything, mock.Anything,
	).Return(&s3.PutObjectOutput{}, uploadError)

	uploader := NewS3Uploader(mockClient, context.TODO())

	result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip")

	assert.Equal(t, uploadError, result.Error)
	mockClient.AssertExpectations(t)
}
