package terraform_test

import (
	"github.com/lewis-od/wavelength/internal/executor"
	"github.com/lewis-od/wavelength/internal/terraform"
	"github.com/lewis-od/wavelength/internal/testutil/mock_executor"
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
	mockExecutor := new(mock_executor.MockExecutor)
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
