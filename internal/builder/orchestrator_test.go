package builder_test

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/testutil/mock_builder"
	"github.com/lewis-od/wavelength/internal/testutil/mock_printer"
	"github.com/lewis-od/wavelength/internal/testutil/mock_uploader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestOrchestrator(t *testing.T) {
	lambdas := []string{"one", "two"}

	var mockBuilder *mock_builder.MockBuilder
	var mockUploader *mock_uploader.MockUploader
	var mockPrinter *mock_printer.MockPrinter
	var mockOrchestrator builder.Orchestrator

	setupTest := func() {
		mockBuilder = new(mock_builder.MockBuilder)
		mockUploader = new(mock_uploader.MockUploader)
		mockPrinter = new(mock_printer.MockPrinter)
		mockPrinter.On("Printlnf", mock.Anything, mock.Anything).Return()
		mockPrinter.On("Println", mock.Anything).Return()
		mockOrchestrator = builder.NewOrchestrator(mockBuilder, mockUploader, mockPrinter)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, mockBuilder, mockUploader, mockPrinter)
	}

	t.Run("BuildLambdas", func(t *testing.T) {
		setupTest()
		successResult := &builder.BuildResult{
			LambdaName: "foo",
			Error:  nil,
			Output: []byte("success"),
		}
		mockBuilder.On("BuildLambda", "one").Return(successResult)
		mockBuilder.On("BuildLambda", "two").Return(successResult)

		err := mockOrchestrator.BuildLambdas(lambdas)

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
	t.Run("UploadLambdas", func(t *testing.T) {
		version := "version"
		bucketName := "bucketName"
		setupTest()

		mockUploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(nil)
		mockUploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(nil)

		err := mockOrchestrator.UploadLambdas(version, bucketName, lambdas)

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
}
