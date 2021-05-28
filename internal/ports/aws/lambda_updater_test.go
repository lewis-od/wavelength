package aws_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/mocks"
	"github.com/lewis-od/wavelength/internal/ports/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockUpdateFunctionCodeAPI struct {
	mock.Mock
}

func (m *mockUpdateFunctionCodeAPI) UpdateFunctionCode(
	ctx context.Context,
	params *lambda.UpdateFunctionCodeInput,
	optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionCodeOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*lambda.UpdateFunctionCodeOutput), args.Error(1)
}

func TestLambdaUpdater_UpdateCode_Success(t *testing.T) {
	lambdaName := "lambda-name"
	bucketName := "bucket-name"
	bucketKey := "key"
	expectedInput := &lambda.UpdateFunctionCodeInput{
		FunctionName: &lambdaName,
		S3Bucket:     &bucketName,
		S3Key:        &bucketKey,
	}

	var mockClient *mockUpdateFunctionCodeAPI
	var mockProviderFactory *mocks.MockAssumeRoleProviderFactory
	var updater builder.Updater

	setupTest := func() {
		mockClient = new(mockUpdateFunctionCodeAPI)
		mockProviderFactory = new(mocks.MockAssumeRoleProviderFactory)
		updater = aws.NewLambdaUpdater(mockClient, mockProviderFactory, context.TODO())
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mockClient.AssertExpectations(t)
		mockProviderFactory.AssertExpectations(t)
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		mockClient.On(
			"UpdateFunctionCode",
			context.TODO(), expectedInput, mock.Anything,
		).Return(&lambda.UpdateFunctionCodeOutput{}, nil)

		err := updater.UpdateCode(lambdaName, bucketName, bucketKey, nil)

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
	t.Run("Error", func(t *testing.T) {
		setupTest()
		expectedErr := fmt.Errorf("error")
		mockClient.On(
			"UpdateFunctionCode",
			context.TODO(), expectedInput, mock.Anything,
		).Return(&lambda.UpdateFunctionCodeOutput{}, expectedErr)

		err := updater.UpdateCode(lambdaName, bucketName, bucketKey, nil)

		assert.Equal(t, expectedErr, err)
		assertExpectationsOnMocks(t)
	})
	t.Run("AssumeRole", func(t *testing.T) {
		setupTest()
		options := &lambda.Options{}
		mockClient.On(
			"UpdateFunctionCode",
			context.TODO(), expectedInput, mock.Anything,
		).Run(func(arguments mock.Arguments) {
			optFns := arguments.Get(2).([]func(*lambda.Options))
			for _, optFn := range optFns {
				optFn(options)
			}
		}).Return(&lambda.UpdateFunctionCodeOutput{}, nil)
		mockProviderFactory.On("CreateProvider", "my-role").Return(&stscreds.AssumeRoleProvider{})

		err := updater.UpdateCode(lambdaName, bucketName, bucketKey, &builder.Role{RoleID: "my-role"})

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
}
