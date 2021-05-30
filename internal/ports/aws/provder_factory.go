package aws

import (
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumeRoleProviderFactory interface {
	CreateProvider(roleArn string) *stscreds.AssumeRoleProvider
}

type assumeRoleProviderFactory struct {
	client *sts.Client
}

func NewAssumeRoleProviderFactory(client *sts.Client) AssumeRoleProviderFactory {
	return &assumeRoleProviderFactory{client: client}
}

func (f *assumeRoleProviderFactory) CreateProvider(roleArn string) *stscreds.AssumeRoleProvider {
	return stscreds.NewAssumeRoleProvider(f.client, roleArn)
}
