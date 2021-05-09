package mocks

import "github.com/stretchr/testify/mock"

type MockFinder struct {
	mock.Mock
}

func (m *MockFinder) FindArtifactBucketName() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockFinder) FindLambdas(lambdaNames []string) ([]string, error) {
	args := m.Called(lambdaNames)
	return args.Get(0).([]string), args.Error(1)
}
