package builder_test

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/mocks"
	"github.com/lewis-od/wavelength/internal/progress"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestOrchestrator(t *testing.T) {
	lambdas := []string{"one", "two"}

	var mockBuilder *mocks.MockBuilder
	var mockUploader *mocks.MockUploader
	var mockDisplay *mocks.MockDisplay
	var mockPrinter *mocks.MockPrinter
	var orchestrator builder.Orchestrator

	setupTest := func() {
		mockBuilder = new(mocks.MockBuilder)
		mockUploader = new(mocks.MockUploader)
		mockDisplay = new(mocks.MockDisplay)
		mockPrinter = new(mocks.MockPrinter)
		mockPrinter.On("Printlnf", mock.Anything, mock.Anything).Return()
		mockPrinter.On("Println", mock.Anything).Return()
		orchestrator = builder.NewOrchestrator(mockBuilder, mockUploader, mockDisplay, mockPrinter)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, mockBuilder, mockUploader, mockDisplay)
	}

	t.Run("BuildLambdas", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			setupTest()
			successOne := &builder.BuildResult{
				LambdaName: "one",
				Error:      nil,
				Output:     []byte("success"),
			}
			successTwo := &builder.BuildResult{
				LambdaName: "two",
				Error:      nil,
				Output:     []byte("success"),
			}
			mockBuilder.On("BuildLambda", "one").Return(successOne)
			mockBuilder.On("BuildLambda", "two").Return(successTwo)

			mockDisplay.On("Init", progress.Build).Return()
			mockDisplay.On("Started", "one").Return()
			mockDisplay.On("Started", "two").Return()
			mockDisplay.On("Completed", "one", true).Return()
			mockDisplay.On("Completed", "two", true).Return()

			failedBuilds := orchestrator.BuildLambdas(lambdas)

			assert.Empty(t, failedBuilds)
			assertExpectationsOnMocks(t)
		})
		t.Run("ErrorBuilding", func(t *testing.T) {
			setupTest()
			successResult := &builder.BuildResult{
				LambdaName: "one",
				Error:      nil,
				Output:     []byte("success"),
			}
			errorResult := &builder.BuildResult{
				LambdaName: "two",
				Error:      fmt.Errorf("error"),
				Output:     []byte("error"),
			}
			mockBuilder.On("BuildLambda", "one").Return(successResult)
			mockBuilder.On("BuildLambda", "two").Return(errorResult)

			mockDisplay.On("Init", progress.Build).Return()
			mockDisplay.On("Started", "one").Return()
			mockDisplay.On("Started", "two").Return()
			mockDisplay.On("Completed", "one", true).Return()
			mockDisplay.On("Completed", "two", false).Return()

			failedBuilds := orchestrator.BuildLambdas(lambdas)

			assert.Len(t, failedBuilds, 1)
			assert.Contains(t, failedBuilds, errorResult)
			assertExpectationsOnMocks(t)
		})
	})
	t.Run("UploadLambdas", func(t *testing.T) {
		version := "version"
		bucketName := "bucketName"

		t.Run("Success", func(t *testing.T) {
			setupTest()

			successOne := &builder.BuildResult{
				LambdaName: "one",
				Error:      nil,
				Output:     []byte(""),
			}
			successTwo := &builder.BuildResult{
				LambdaName: "two",
				Error:      nil,
				Output:     []byte(""),
			}
			mockUploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(successOne)
			mockUploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(successTwo)

			mockDisplay.On("Init", progress.Upload).Return()
			mockDisplay.On("Started", "one").Return()
			mockDisplay.On("Started", "two").Return()
			mockDisplay.On("Completed", "one", true).Return()
			mockDisplay.On("Completed", "two", true).Return()

			failedUploads := orchestrator.UploadLambdas(version, bucketName, lambdas)

			assert.Empty(t, failedUploads)
			assertExpectationsOnMocks(t)
		})
		t.Run("Error", func(t *testing.T) {
			setupTest()

			successResult := &builder.BuildResult{
				LambdaName: "one",
				Error:      nil,
				Output:     []byte(""),
			}
			errorResult := &builder.BuildResult{
				LambdaName: "two",
				Error:      fmt.Errorf("error uploading"),
				Output:     []byte(""),
			}
			mockUploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(successResult)
			mockUploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(errorResult)

			mockDisplay.On("Init", progress.Upload).Return()
			mockDisplay.On("Started", "one").Return()
			mockDisplay.On("Started", "two").Return()
			mockDisplay.On("Completed", "one", true).Return()
			mockDisplay.On("Completed", "two", false).Return()

			failedUploads := orchestrator.UploadLambdas(version, bucketName, lambdas)

			assert.Contains(t, failedUploads, errorResult)
			assertExpectationsOnMocks(t)
		})
	})
}
