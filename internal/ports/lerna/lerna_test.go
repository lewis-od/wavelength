package lerna_test

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/executor"
	"github.com/lewis-od/wavelength/internal/mocks"
	"github.com/lewis-od/wavelength/internal/ports/lerna"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildLambda(t *testing.T) {
	projectName := "project"
	lambdaName := "lambda"

	expectedContext := &executor.CommandContext{Directory: "."}

	buildOutput := []byte("output")
	mockExecutor := new(mocks.MockExecutor)
	mockExecutor.On(
		"ExecuteAndCapture",
		[]string{"run", "build", "--scope", "@project/lambda", "--include-dependencies"},
		expectedContext,
	).Return(buildOutput, nil)
	l := lerna.NewLerna(mockExecutor, &projectName)

	result := l.BuildLambda(lambdaName)

	expectedResult := &builder.BuildResult{
		LambdaName: lambdaName,
		Error:  nil,
		Output: buildOutput,
	}
	assert.Equal(t, expectedResult, result)
	mockExecutor.AssertExpectations(t)
}
