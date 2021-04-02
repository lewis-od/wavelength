package terraform

import (
	"github.com/lewis-od/lambda-build/pkg/executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockExecutor struct {
	mock.Mock
}

func (e *MockExecutor) ExecuteWithContext(arguments []string, context *executor.CommandContext) ([]byte, error) {
	args := e.Called(arguments, context)
	return args.Get(0).([]byte), args.Error(1)
}

const dummyOutput string = `
{
	"my output": {
		"sensitive": false,
		"type": "string",
		"value": "output value"
	}
}
`

func TestOutput(t *testing.T) {
	directoryName := "directory"
	mockExecutor := new(MockExecutor)
	mockExecutor.On(
		"ExecuteWithContext",
		[]string{"output", "-json"},
		&executor.CommandContext{Directory: directoryName},
	).Return([]byte(dummyOutput), nil)

	tf := NewTerraform(mockExecutor)
	outputs, err := tf.output(directoryName)

	assert.Nil(t, err)
	assert.Equal(t, outputs["my output"].Value, "output value")
}
