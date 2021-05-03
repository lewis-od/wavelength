package service

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/testutil/mock_finder"
	"github.com/lewis-od/wavelength/internal/testutil/mock_printer"
	"github.com/lewis-od/wavelength/internal/testutil/mock_updater"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUpdateService_Run(t *testing.T) {
	projectName := "my-project"
	version := "version"
	lambdaOne := "one"
	lambdaTwo := "two"
	lambdas := []string{lambdaOne, lambdaTwo}
	bucketName := "my-bucket"

	var finder *mock_finder.MockFinder
	var updater *mock_updater.MockUpdater
	var printer *mock_printer.MockPrinter
	var updateService UpdateService

	setupTest := func() {
		finder = new(mock_finder.MockFinder)
		updater = new(mock_updater.MockUpdater)
		printer = new(mock_printer.MockPrinter)
		updateService = NewUpdateService(finder, updater, printer, &projectName)
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On(
			"Printlnf",
			"⬆️ Updating %s with code at s3://%s/%s", []interface{}{lambdaOne, bucketName, "version/one.zip"},
		).Return()
		printer.On(
			"Printlnf",
			"⬆️ Updating %s with code at s3://%s/%s", []interface{}{lambdaTwo, bucketName, "version/two.zip"},
		).Return()
		printer.On("Printlnf", "✅ Successfully updated %s", []interface{}{lambdaOne}).Return()
		printer.On("Printlnf", "✅ Successfully updated %s", []interface{}{lambdaTwo}).Return()

		updater.On("UpdateCode", lambdaOne, bucketName, "version/one.zip").Return(nil)
		updater.On("UpdateCode", lambdaTwo, bucketName, "version/two.zip").Return(nil)

		updateService.Run(version, lambdas)

		mock.AssertExpectationsForObjects(t, finder, updater, printer)
	})
	t.Run("UploadError", func(t *testing.T) {
		setupTest()
		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		uploadErr := fmt.Errorf("error updating lambda")
		updater.On("UpdateCode", lambdaOne, bucketName, "version/one.zip").Return(uploadErr)

		printer.On(
			"Printlnf",
			"⬆️ Updating %s with code at s3://%s/%s", []interface{}{lambdaOne, bucketName, "version/one.zip"},
		).Return()
		printer.On("PrintErr", uploadErr).Return()

		updateService.Run(version, lambdas)

		mock.AssertExpectationsForObjects(t, finder, updater, printer)
	})
}
