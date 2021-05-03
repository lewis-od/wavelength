package mock_builder

import "github.com/stretchr/testify/mock"

type MockBuilder struct {
	mock.Mock
}

func (m *MockBuilder) BuildLambda(lambdaName string) error {
	args := m.Called(lambdaName)
	return args.Error(0)
}
