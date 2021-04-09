package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/lewis-od/lambda-build/pkg/builder"
)

type UpdateFunctionCodeApi interface {
	UpdateFunctionCode(ctx context.Context,
		params *lambda.UpdateFunctionCodeInput,
		optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionCodeOutput, error)
}

type lambdaUpdater struct {
	client UpdateFunctionCodeApi
	ctx    context.Context
}

func NewLambdaUpdater(client UpdateFunctionCodeApi, ctx context.Context) builder.Updater {
	return &lambdaUpdater{
		client: client,
		ctx:    ctx,
	}
}

func (l *lambdaUpdater) UpdateCode(lambdaName, bucketName, bucketKey string) error {
	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: &lambdaName,
		S3Bucket:     &bucketName,
		S3Key:        &bucketKey,
	}
	_, err := l.client.UpdateFunctionCode(l.ctx, input)
	return err
}
