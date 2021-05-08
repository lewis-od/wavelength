package lerna

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/executor"
	"github.com/lewis-od/wavelength/internal/testutil/mock_executor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildLambda(t *testing.T) {
	projectName := "project"
	lambdaName := "lambda"

	expectedContext := &executor.CommandContext{Directory: "."}

	buildOutput := []byte("output")
	mockExecutor := new(mock_executor.MockExecutor)
	mockExecutor.On(
		"ExecuteAndCapture",
		[]string{"run", "build", "--scope", "@project/lambda", "--include-dependencies"},
		expectedContext,
	).Return(buildOutput, nil)
	lerna := NewLerna(mockExecutor, &projectName)

	result := lerna.BuildLambda(lambdaName)

	expectedResult := &builder.BuildResult{
		LambdaName: lambdaName,
		Error:  nil,
		Output: buildOutput,
	}
	assert.Equal(t, expectedResult, result)
	mockExecutor.AssertExpectations(t)
}
