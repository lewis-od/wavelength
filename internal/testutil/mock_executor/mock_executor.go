package mock_executor

import (
	"github.com/lewis-od/wavelength/internal/executor"
	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

func (e *MockExecutor) Execute(arguments []string) error {
	args := e.Called(arguments)
	return args.Error(0)
}

func (e *MockExecutor) ExecuteAndCapture(arguments []string, context *executor.CommandContext) ([]byte, error) {
	args := e.Called(arguments, context)
	return args.Get(0).([]byte), args.Error(1)
}
