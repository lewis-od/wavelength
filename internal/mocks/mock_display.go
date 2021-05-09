package mocks

import (
	"github.com/lewis-od/wavelength/internal/progress"
	"github.com/stretchr/testify/mock"
)

type MockDisplay struct {
	mock.Mock
}

func (m *MockDisplay) Init(action progress.Action) {
	m.Called(action)
}

func (m *MockDisplay) Started(lambdaName string) {
	m.Called(lambdaName)
}

func (m *MockDisplay) Completed(lambdaName string, wasSuccessful bool) {
	m.Called(lambdaName, wasSuccessful)
}
