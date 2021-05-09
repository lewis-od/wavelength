package mock_display

import "github.com/stretchr/testify/mock"

type MockDisplay struct {
	mock.Mock
}

func (m *MockDisplay) Started(lambdaName string) {
	m.Called(lambdaName)
}

func (m *MockDisplay) Completed(lambdaName string, wasSuccessful bool) {
	m.Called(lambdaName, wasSuccessful)
}

