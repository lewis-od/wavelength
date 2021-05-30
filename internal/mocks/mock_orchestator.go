package mocks

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/stretchr/testify/mock"
)

type MockOrchestrator struct {
	mock.Mock
}

func (m *MockOrchestrator) BuildLambdas(lambdas []string) []*builder.BuildResult {
	args := m.Called(lambdas)
	return args.Get(0).([]*builder.BuildResult)
}

func (m *MockOrchestrator) UploadLambdas(version, bucketName string, lambdas []string, role *builder.Role) []*builder.BuildResult {
	args := m.Called(version, bucketName, lambdas, role)
	return args.Get(0).([]*builder.BuildResult)
}
