package mock_orchestrator

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

func (m *MockOrchestrator) UploadLambdas(version, bucketName string, lambdas []string) []*builder.BuildResult {
	args := m.Called(version, bucketName, lambdas)
	return args.Get(0).([]*builder.BuildResult)
}
