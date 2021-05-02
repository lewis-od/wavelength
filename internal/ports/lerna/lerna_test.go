package lerna

import (
	"github.com/lewis-od/lambda-build/internal/executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockExecutor struct {
	mock.Mock
}

func (e *MockExecutor) Execute(arguments []string) (error) {
	args := e.Called(arguments)
	return args.Error(0)
}

func (e *MockExecutor) ExecuteAndCapture(arguments []string, context *executor.CommandContext) ([]byte, error) {
	args := e.Called(arguments, context)
	return args.Get(0).([]byte), args.Error(1)
}

func TestBuildLambda(t *testing.T) {
	mockExecutor := new(MockExecutor)
	mockExecutor.On(
		"Execute",
		[]string{"run", "build", "--scope", "@project/lambda", "--include-dependencies",
	}).Return(nil)
	lerna := NewLerna(mockExecutor, "project")

	err := lerna.BuildLambda("lambda")

	assert.Nil(t, err)
}
