package mocks

import "github.com/stretchr/testify/mock"

type MockUpdater struct {
	mock.Mock
}

func (m *MockUpdater) UpdateCode(lambdaName, bucketName, bucketKey string) error {
	args := m.Called(lambdaName, bucketName, bucketKey)
	return args.Error(0)
}
