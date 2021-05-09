package aws_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/ports/aws"
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

func TestS3Uploader_UploadLambda(t *testing.T) {
	var mockClient *mockS3Client
	var uploader builder.Uploader

	setupTest := func() {
		mockClient = new(mockS3Client)
		uploader = aws.NewS3Uploader(mockClient, context.TODO())
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		mockClient.On(
			"PutObject",
			context.TODO(), mock.Anything, mock.Anything, mock.Anything,
		).Return(&s3.PutObjectOutput{}, nil)

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip")

		assert.Nil(t, result.Error)
		mockClient.AssertExpectations(t)
		assert.Equal(t, "test/my-lambda.zip", *mockClient.input.Key)
	})
	t.Run("FileOpenError", func(t *testing.T) {
		setupTest()

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/missing.zip")

		assert.NotNil(t, result.Error)
		mockClient.AssertExpectations(t)
	})
	t.Run("D3PutError", func(t *testing.T) {
		setupTest()
		uploadError := fmt.Errorf("upload error")
		mockClient.On(
			"PutObject",
			context.TODO(), mock.Anything, mock.Anything, mock.Anything,
		).Return(&s3.PutObjectOutput{}, uploadError)

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip")

		assert.Equal(t, uploadError, result.Error)
		mockClient.AssertExpectations(t)
	})
}
