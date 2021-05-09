package mock_builder

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/stretchr/testify/mock"
)

type MockBuilder struct {
	mock.Mock
}

func (m *MockBuilder) BuildLambda(lambdaName string) *builder.BuildResult {
	args := m.Called(lambdaName)
	return args.Get(0).(*builder.BuildResult)
}
