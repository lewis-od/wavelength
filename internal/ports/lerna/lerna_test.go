package lerna

import (
	"github.com/lewis-od/lambda-build/internal/testutil/mock_executor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildLambda(t *testing.T) {
	mockExecutor := new(mock_executor.MockExecutor)
	mockExecutor.On(
		"Execute",
		[]string{"run", "build", "--scope", "@project/lambda", "--include-dependencies"},
	).Return(nil)
	lerna := NewLerna(mockExecutor, "project")

	err := lerna.BuildLambda("lambda")

	assert.Nil(t, err)
	mockExecutor.AssertExpectations(t)
}
