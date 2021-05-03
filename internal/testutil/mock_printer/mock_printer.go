package mock_printer

import "github.com/stretchr/testify/mock"

type MockPrinter struct {
	mock.Mock
}

func (n *MockPrinter) Println(a ...interface{}) {
	n.Called(a)
}

func (n *MockPrinter) Printlnf(format string, a ...interface{}) {
	n.Called(format, a)
}

func (n *MockPrinter) PrintErr(err error) {
	n.Called(err)
}
