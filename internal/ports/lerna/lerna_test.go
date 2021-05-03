package lerna

import (
	"github.com/lewis-od/wavelength/internal/testutil/mock_executor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildLambda(t *testing.T) {
	projectName := "project"
	lambdaName := "lambda"

	mockExecutor := new(mock_executor.MockExecutor)
	mockExecutor.On(
		"Execute",
		[]string{"run", "build", "--scope", "@project/lambda", "--include-dependencies"},
	).Return(nil)
	lerna := NewLerna(mockExecutor, &projectName)

	err := lerna.BuildLambda(lambdaName)

	assert.Nil(t, err)
	mockExecutor.AssertExpectations(t)
}
