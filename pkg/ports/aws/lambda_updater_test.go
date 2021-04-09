package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
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

	mockClient := new(mockUpdateFunctionCodeAPI)
	mockClient.On(
		"UpdateFunctionCode",
		context.TODO(), expectedInput, mock.Anything,
	).Return(&lambda.UpdateFunctionCodeOutput{}, nil)

	updater := NewLambdaUpdater(mockClient, context.TODO())

	err := updater.UpdateCode(lambdaName, bucketName, bucketKey)

	assert.Nil(t, err)
	mockClient.AssertExpectations(t)
}

func TestLambdaUpdater_UpdateCode_Error(t *testing.T) {
	lambdaName := "lambda-name"
	bucketName := "bucket-name"
	bucketKey := "key"
	expectedInput := &lambda.UpdateFunctionCodeInput{
		FunctionName: &lambdaName,
		S3Bucket:     &bucketName,
		S3Key:        &bucketKey,
	}

	mockClient := new(mockUpdateFunctionCodeAPI)
	expectedErr := fmt.Errorf("error")
	mockClient.On(
		"UpdateFunctionCode",
		context.TODO(), expectedInput, mock.Anything,
	).Return(&lambda.UpdateFunctionCodeOutput{}, expectedErr)

	updater := NewLambdaUpdater(mockClient, context.TODO())

	err := updater.UpdateCode(lambdaName, bucketName, bucketKey)

	assert.Equal(t, expectedErr, err)
	mockClient.AssertExpectations(t)
}
