package service_test

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/mocks"
	"github.com/lewis-od/wavelength/internal/service"
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
	roleToAssume := &builder.Role{RoleID: "my-role"}

	var finder *mocks.MockFinder
	var updater *mocks.MockUpdater
	var printer *mocks.MockPrinter
	var updateService service.UpdateService

	setupTest := func() {
		finder = new(mocks.MockFinder)
		updater = new(mocks.MockUpdater)
		printer = new(mocks.MockPrinter)
		updateService = service.NewUpdateService(finder, updater, printer, &projectName)
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

		updater.On("UpdateCode", lambdaOne, bucketName, "version/one.zip", roleToAssume).Return(nil)
		updater.On("UpdateCode", lambdaTwo, bucketName, "version/two.zip", roleToAssume).Return(nil)

		updateService.Run(version, lambdas, roleToAssume)

		mock.AssertExpectationsForObjects(t, finder, updater, printer)
	})
	t.Run("UploadError", func(t *testing.T) {
		setupTest()
		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		uploadErr := fmt.Errorf("error updating lambda")
		updater.On("UpdateCode", lambdaOne, bucketName, "version/one.zip", roleToAssume).Return(uploadErr)

		printer.On(
			"Printlnf",
			"⬆️ Updating %s with code at s3://%s/%s", []interface{}{lambdaOne, bucketName, "version/one.zip"},
		).Return()
		printer.On("PrintErr", uploadErr).Return()

		updateService.Run(version, lambdas, roleToAssume)

		mock.AssertExpectationsForObjects(t, finder, updater, printer)
	})
}
