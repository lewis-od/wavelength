package mock_error_logger

import (
	"github.com/lewis-od/wavelength/internal/error_logger"
	"github.com/stretchr/testify/mock"
)

type MockErrorLogger struct {
	mock.Mock
}

func (m *MockErrorLogger) AddError(wavelengthError *error_logger.WavelengthError) {
	m.Called(wavelengthError)
}

func (m *MockErrorLogger) WriteLogFile() error {
	args := m.Called()
	return args.Error(0)
}
