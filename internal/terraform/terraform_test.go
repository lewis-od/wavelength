package terraform_test

import (
	"github.com/lewis-od/wavelength/internal/executor"
	"github.com/lewis-od/wavelength/internal/mocks"
	"github.com/lewis-od/wavelength/internal/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	mockExecutor := new(mocks.MockExecutor)
	mockExecutor.On(
		"ExecuteAndCapture",
		[]string{"output", "-json"},
		&executor.CommandContext{Directory: directoryName},
	).Return([]byte(dummyOutput), nil)

	tf := terraform.NewTerraform(mockExecutor)
	outputs, err := tf.Output(directoryName)

	assert.Nil(t, err)
	assert.Equal(t, outputs["my output"].Value, "output value")
}
