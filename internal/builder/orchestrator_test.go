package builder

import (
	"github.com/lewis-od/wavelength/internal/testutil/mock_builder"
	"github.com/lewis-od/wavelength/internal/testutil/mock_printer"
	"github.com/lewis-od/wavelength/internal/testutil/mock_uploader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestOrchestrator(t *testing.T) {
	lambdas := []string{"one", "two"}

	var builder *mock_builder.MockBuilder
	var uploader *mock_uploader.MockUploader
	var printer *mock_printer.MockPrinter
	var orchestrator Orchestrator

	setupTest := func() {
		builder = new(mock_builder.MockBuilder)
		uploader = new(mock_uploader.MockUploader)
		printer = new(mock_printer.MockPrinter)
		printer.On("Printlnf", mock.Anything, mock.Anything).Return()
		printer.On("Println", mock.Anything).Return()
		orchestrator = NewOrchestrator(builder, uploader, printer)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, builder, uploader, printer)
	}

	t.Run("BuildLambdas", func(t *testing.T) {
		setupTest()
		builder.On("BuildLambda", "one").Return(nil)
		builder.On("BuildLambda", "two").Return(nil)

		err := orchestrator.BuildLambdas(lambdas)

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
	t.Run("UploadLambdas", func(t *testing.T) {
		version := "version"
		bucketName := "bucketName"
		setupTest()

		uploader.On("UploadLambda", version, bucketName, "one", "lambdas/one/dist/one.zip").Return(nil)
		uploader.On("UploadLambda", version, bucketName, "two", "lambdas/two/dist/two.zip").Return(nil)

		err := orchestrator.UploadLambdas(version, bucketName, lambdas)

		assert.Nil(t, err)
		assertExpectationsOnMocks(t)
	})
}
