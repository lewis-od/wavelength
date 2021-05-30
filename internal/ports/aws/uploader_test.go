package aws_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/mocks"
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
	var s3Client *mockS3Client
	var roleProviderFactory *mocks.MockAssumeRoleProviderFactory
	var uploader builder.Uploader

	setupTest := func() {
		s3Client = new(mockS3Client)
		roleProviderFactory = new(mocks.MockAssumeRoleProviderFactory)
		uploader = aws.NewS3Uploader(s3Client, roleProviderFactory, context.TODO())
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, s3Client, roleProviderFactory)
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		s3Client.On(
			"PutObject",
			context.TODO(), mock.Anything, mock.Anything, mock.Anything,
		).Return(&s3.PutObjectOutput{}, nil)

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip", nil)

		assert.Nil(t, result.Error)
		assertExpectationsOnMocks(t)
		assert.Equal(t, "test/my-lambda.zip", *s3Client.input.Key)
	})
	t.Run("FileOpenError", func(t *testing.T) {
		setupTest()

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/missing.zip", nil)

		assert.NotNil(t, result.Error)
		assertExpectationsOnMocks(t)
	})
	t.Run("S3PutError", func(t *testing.T) {
		setupTest()
		uploadError := fmt.Errorf("upload error")
		s3Client.On(
			"PutObject",
			context.TODO(), mock.Anything, mock.Anything, mock.Anything,
		).Return(&s3.PutObjectOutput{}, uploadError)

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip", nil)

		assert.Equal(t, uploadError, result.Error)
		assertExpectationsOnMocks(t)
	})
	t.Run("AssumeRole", func(t *testing.T) {
		setupTest()
		roleToAssume := &builder.Role{RoleID: "my-role"}
		options := &s3.Options{}
		s3Client.On(
			"PutObject",
			context.TODO(), mock.Anything, mock.Anything, mock.Anything,
		).Run(func(arguments mock.Arguments) {
			optFns := arguments.Get(2).([]func(options *s3.Options))
			for _, optFn := range optFns {
				optFn(options)
			}
		}).Return(&s3.PutObjectOutput{}, nil)
		roleProvider := &stscreds.AssumeRoleProvider{}
		roleProviderFactory.On("CreateProvider", "my-role").Return(roleProvider)

		result := uploader.UploadLambda("test", "artifact-bucket", "my-lambda", "testdata/artifact.zip", roleToAssume)

		assert.Nil(t, result.Error)
		assertExpectationsOnMocks(t)
		assert.Equal(t, roleProvider, options.Credentials)
		assert.Equal(t, "test/my-lambda.zip", *s3Client.input.Key)
	})
}
