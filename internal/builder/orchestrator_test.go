package builder

import (
	"github.com/lewis-od/lambda-build/internal/testutil/mock_builder"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_printer"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_uploader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var lambdas = []string{"one", "two"}
var version = "version"
var bucketName = "bucketName"

var builder *mock_builder.MockBuilder
var uploader *mock_uploader.MockUploader
var printer *mock_printer.MockPrinter

func initMocks() {
	builder = new(mock_builder.MockBuilder)
	uploader = new(mock_uploader.MockUploader)
	printer = new(mock_printer.MockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Println", mock.Anything).Return()
}

func assertExpectationsOnMocks(t *testing.T) {
	mock.AssertExpectationsForObjects(t, builder, uploader, printer)
}

func TestBuildLambdas_Success(t *testing.T) {
	initMocks()
	builder.On("BuildLambda", "one").Return(nil)
	builder.On("BuildLambda", "two").Return(nil)

	orchestrator := NewOrchestrator(builder, uploader, printer)
	err := orchestrator.BuildLambdas(lambdas)

	assert.Nil(t, err)
	assertExpectationsOnMocks(t)
}

func TestUploadLambdas_Success(t *testing.T) {
	initMocks()
	uploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(nil)
	uploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(nil)

	orchestrator := NewOrchestrator(builder, uploader, printer)
	err := orchestrator.UploadLambdas(version, bucketName, lambdas)

	assert.Nil(t, err)
	assertExpectationsOnMocks(t)
}
