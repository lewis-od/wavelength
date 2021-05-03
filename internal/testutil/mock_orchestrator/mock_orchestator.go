package mock_orchestrator

import "github.com/stretchr/testify/mock"

type MockOrchestrator struct {
	mock.Mock
}

func (m *MockOrchestrator) BuildLambdas(lambdas []string) error {
	args := m.Called(lambdas)
	return args.Error(0)
}

func (m *MockOrchestrator) UploadLambdas(version, bucketName string, lambdas []string) error {
	args := m.Called(version, bucketName, lambdas)
	return args.Error(0)
}
