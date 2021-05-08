package mock_builder

import "github.com/stretchr/testify/mock"

type MockBuilder struct {
	mock.Mock
}

func (m *MockBuilder) BuildLambda(lambdaName string) ([]byte, error) {
	args := m.Called(lambdaName)
	return args.Get(0).([]byte), args.Error(1)
}
