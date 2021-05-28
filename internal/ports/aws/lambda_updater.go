package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/lewis-od/wavelength/internal/builder"
)

type UpdateFunctionCodeApi interface {
	UpdateFunctionCode(ctx context.Context,
		params *lambda.UpdateFunctionCodeInput,
		optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionCodeOutput, error)
}

type lambdaUpdater struct {
	client          UpdateFunctionCodeApi
	providerFactory AssumeRoleProviderFactory
	ctx             context.Context
}

func NewLambdaUpdater(
	client UpdateFunctionCodeApi,
	providerFactory AssumeRoleProviderFactory,
	ctx context.Context,
) builder.Updater {
	return &lambdaUpdater{
		client:          client,
		providerFactory: providerFactory,
		ctx:             ctx,
	}
}

func (l *lambdaUpdater) UpdateCode(lambdaName, bucketName, bucketKey string, role *builder.AssumeRole) error {
	input := &lambda.UpdateFunctionCodeInput{
		FunctionName: &lambdaName,
		S3Bucket:     &bucketName,
		S3Key:        &bucketKey,
	}

	optFun := func(options *lambda.Options) {}
	if role != nil {
		optFun = func(options *lambda.Options) {
			options.Credentials = l.providerFactory.CreateProvider(role.RoleID)
		}
	}

	_, err := l.client.UpdateFunctionCode(l.ctx, input, optFun)
	return err
}
