package mock_terraform

import (
	"github.com/lewis-od/lambda-build/internal/terraform"
	"github.com/stretchr/testify/mock"
)

type MockTerraform struct {
	mock.Mock
}

func (m *MockTerraform) Output(directory string) (map[string]terraform.Output, error) {
	args := m.Called(directory)
	return args.Get(0).(map[string]terraform.Output), args.Error(1)
}
