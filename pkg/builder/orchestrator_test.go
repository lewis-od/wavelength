package builder

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockBuilder struct {
	mock.Mock
}

func (m *mockBuilder) BuildLambda(lambdaName string) error {
	args := m.Called(lambdaName)
	return args.Error(0)
}

type mockUploader struct {
	mock.Mock
}

func (m *mockUploader) UploadLambda(version, bucketName, lambdaName, artifactLocation string) error {
	args := m.Called(version, bucketName, lambdaName, artifactLocation)
	return args.Error(0)
}

type mockPrinter struct {
	mock.Mock
}

func (n *mockPrinter) Println(a ...interface{}) {
	n.Called(a)
}

func (n *mockPrinter) Printlnf(format string, a ...interface{}) {
	n.Called(format, a)
}

func (n *mockPrinter) PrintErr(err error) {
	n.Called(err)
}

func newMockPrinter() (printer *mockPrinter) {
	printer = new(mockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Println", mock.Anything).Return()
	return
}

func TestBuildLambdas_Success(t *testing.T) {
	lambdas := []string{"one", "two"}

	builder := new(mockBuilder)
	builder.On("BuildLambda", "one").Return(nil)
	builder.On("BuildLambda", "two").Return(nil)

	uploader := new(mockUploader)
	printer := newMockPrinter()

	orchestrator := NewOrchestrator(builder, uploader, printer)
	err := orchestrator.BuildLambdas(lambdas)

	assert.Nil(t, err)
	builder.AssertExpectations(t)
	uploader.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestUploadLambdas_Success(t *testing.T) {
	version := "version"
	bucketName := "bucket"
	lambdas := []string{"one", "two"}

	builder := new(mockBuilder)

	uploader := new(mockUploader)
	uploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(nil)
	uploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(nil)

	printer := newMockPrinter()

	orchestrator := NewOrchestrator(builder, uploader, printer)
	err := orchestrator.UploadLambdas(version, bucketName, lambdas)

	assert.Nil(t, err)
	builder.AssertExpectations(t)
	uploader.AssertExpectations(t)
	printer.AssertExpectations(t)
}
