package mocks

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/stretchr/testify/mock"
)

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) UpdateCode(lambdaName, bucketName, bucketKey string, role *builder.AssumeRole) error {
	args := m.Called(lambdaName, bucketName, bucketKey, role)
	return args.Error(0)
}
