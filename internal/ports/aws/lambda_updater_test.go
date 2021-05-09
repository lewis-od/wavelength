package aws_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/ports/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockUpdateFunctionCodeAPI struct {
	mock.Mock
}

func (m *mockUpdateFunctionCodeAPI) UpdateFunctionCode(ctx context.Context,
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
	var updater builder.Updater

	setupTest := func() {
		mockClient = new(mockUpdateFunctionCodeAPI)
		updater = aws.NewLambdaUpdater(mockClient, context.TODO())
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		mockClient.On(
			"UpdateFunctionCode",
			context.TODO(), expectedInput, mock.Anything,
		).Return(&lambda.UpdateFunctionCodeOutput{}, nil)

		err := updater.UpdateCode(lambdaName, bucketName, bucketKey)

		assert.Nil(t, err)
		mockClient.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {
		setupTest()
		expectedErr := fmt.Errorf("error")
		mockClient.On(
			"UpdateFunctionCode",
			context.TODO(), expectedInput, mock.Anything,
		).Return(&lambda.UpdateFunctionCodeOutput{}, expectedErr)

		err := updater.UpdateCode(lambdaName, bucketName, bucketKey)

		assert.Equal(t, expectedErr, err)
		mockClient.AssertExpectations(t)
	})
}
