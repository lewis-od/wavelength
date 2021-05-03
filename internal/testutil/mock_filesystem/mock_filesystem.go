package mock_filesystem

import (
	"github.com/lewis-od/wavelength/internal/io"
	"github.com/stretchr/testify/mock"
)

type MockFilesystem struct {
	mock.Mock
}

func (m *MockFilesystem) ReadDir(dirname string) ([]io.FileInfo, error) {
	args := m.Called(dirname)
	return args.Get(0).([]io.FileInfo), args.Error(1)
}

func (m *MockFilesystem) FileExists(filename string) bool {
	args := m.Called(filename)
	return args.Bool(0)
}
