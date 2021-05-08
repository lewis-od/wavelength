package builder_test

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/error_logger"
	"github.com/lewis-od/wavelength/internal/testutil/mock_builder"
	"github.com/lewis-od/wavelength/internal/testutil/mock_error_logger"
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
	var mockErrorLogger *mock_error_logger.MockErrorLogger
	var mockPrinter *mock_printer.MockPrinter
	var orchestrator builder.Orchestrator

	setupTest := func() {
		mockBuilder = new(mock_builder.MockBuilder)
		mockUploader = new(mock_uploader.MockUploader)
		mockPrinter = new(mock_printer.MockPrinter)
		mockErrorLogger = new(mock_error_logger.MockErrorLogger)
		mockPrinter.On("Printlnf", mock.Anything, mock.Anything).Return()
		mockPrinter.On("Println", mock.Anything).Return()
		orchestrator = builder.NewOrchestrator(mockBuilder, mockUploader, mockErrorLogger, mockPrinter)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, mockBuilder, mockUploader, mockErrorLogger)
	}

	t.Run("BuildLambdas", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			setupTest()
			successResult := &builder.BuildResult{
				LambdaName: "foo",
				Error:  nil,
				Output: []byte("success"),
			}
			mockBuilder.On("BuildLambda", "one").Return(successResult)
			mockBuilder.On("BuildLambda", "two").Return(successResult)

			err := orchestrator.BuildLambdas(lambdas)

			assert.Nil(t, err)
			assertExpectationsOnMocks(t)
		})
		t.Run("ErrorBuilding", func(t *testing.T) {
			setupTest()
			successResult := &builder.BuildResult{
				LambdaName: "one",
				Error:  nil,
				Output: []byte("success"),
			}
			errorResult := &builder.BuildResult{
				LambdaName: "two",
				Error: fmt.Errorf("error"),
				Output: []byte("error"),
			}
			mockBuilder.On("BuildLambda", "one").Return(successResult)
			mockBuilder.On("BuildLambda", "two").Return(errorResult)

			expectedError := &error_logger.WavelengthError{
				Lambda: "two",
				Output: []byte("error"),
			}
			mockErrorLogger.On("AddError", expectedError).Return()
			mockErrorLogger.On("WriteLogFile").Return(nil)

			err := orchestrator.BuildLambdas(lambdas)

			assert.Equal(t, fmt.Errorf("Error building lambdas [two]"), err)
			assertExpectationsOnMocks(t)
		})
		t.Run("ErrorWritingLogFile", func(t *testing.T) {
			setupTest()
			successResult := &builder.BuildResult{
				LambdaName: "one",
				Error:  nil,
				Output: []byte("success"),
			}
			errorResult := &builder.BuildResult{
				LambdaName: "two",
				Error: fmt.Errorf("error"),
				Output: []byte("error"),
			}
			mockBuilder.On("BuildLambda", "one").Return(successResult)
			mockBuilder.On("BuildLambda", "two").Return(errorResult)

			buildError := &error_logger.WavelengthError{
				Lambda: "two",
				Output: []byte("error"),
			}
			mockErrorLogger.On("AddError", buildError).Return()
			expectedErr := fmt.Errorf("error writing log")
			mockErrorLogger.On("WriteLogFile").Return(expectedErr)

			err := orchestrator.BuildLambdas(lambdas)

			assert.Equal(t, expectedErr, err)
			assertExpectationsOnMocks(t)
		})
	})
	t.Run("UploadLambdas", func(t *testing.T) {
		version := "version"
		bucketName := "bucketName"
		setupTest()

		mockUploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(nil)
		mockUploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(nil)

		err := orchestrator.UploadLambdas(version, bucketName, lambdas)

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
}
