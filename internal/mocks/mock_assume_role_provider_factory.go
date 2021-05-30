package mocks

import (
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/stretchr/testify/mock"
)

type MockAssumeRoleProviderFactory struct {
	mock.Mock
}

func (m *MockAssumeRoleProviderFactory) CreateProvider(roleArn string) *stscreds.AssumeRoleProvider {
	args := m.Called(roleArn)
	return args.Get(0).(*stscreds.AssumeRoleProvider)
}
