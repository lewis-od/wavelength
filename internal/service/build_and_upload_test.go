package service_test

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/lewis-od/wavelength/internal/mocks"
	"github.com/lewis-od/wavelength/internal/service"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestBuildAndUploadService_Run(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	bucketName := "some-bucket"
	role := &builder.Role{RoleID: "my-role"}

	var orchestrator *mocks.MockOrchestrator
	var finder *mocks.MockFinder
	var printer *mocks.MockPrinter
	var command service.BuildAndUploadService

	setupTest := func() {
		orchestrator = new(mocks.MockOrchestrator)
		finder = new(mocks.MockFinder)
		printer = new(mocks.MockPrinter)
		command = service.NewBuildAndUploadService(orchestrator, finder, printer)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, orchestrator, finder, printer)
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		orchestrator.On(
			"BuildLambdas",
			lambdas,
		).Return(make([]*builder.BuildResult, 0, 0))
		orchestrator.On(
			"UploadLambdas",
			version, bucketName, lambdas, role,
		).Return(make([]*builder.BuildResult, 0, 0))

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On("Printlnf", mock.Anything, mock.Anything).Return()
		printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()
		printer.On("Println", mock.Anything).Return()

		command.Run(version, lambdas, false, role)

		assertExpectationsOnMocks(t)
	})
	t.Run("SkipBuild", func(t *testing.T) {
		setupTest()
		orchestrator.On(
			"UploadLambdas",
			version, bucketName, lambdas, role,
		).Return(make([]*builder.BuildResult, 0, 0))

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On("Printlnf", mock.Anything, mock.Anything).Return()
		printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()
		printer.On("Println", mock.Anything).Return()

		command.Run(version, lambdas, true, role)

		assertExpectationsOnMocks(t)
	})
	t.Run("BuildError", func(t *testing.T) {
		setupTest()
		buildErr := &builder.BuildResult{
			LambdaName: "lambda-one",
			Error:      fmt.Errorf("error"),
			Output:     []byte("output"),
		}
		orchestrator.On("BuildLambdas", lambdas).Return([]*builder.BuildResult{buildErr})

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
		printer.On("Printlnf", "🪣 Found artifact bucket %s", []interface{}{bucketName}).Return()
		errToPrint := fmt.Errorf("Error building lambda %s\n%s\n", "lambda-one", "output")
		printer.On("PrintErr", errToPrint).Return()

		command.Run(version, lambdas, false, role)

		assertExpectationsOnMocks(t)
	})
	t.Run("FindArtifactBucketError", func(t *testing.T) {
		setupTest()

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finderErr := fmt.Errorf("error finding artifact bucket")
		finder.On("FindArtifactBucketName").Return(bucketName, finderErr)

		printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
		printer.On("PrintErr", finderErr).Return()

		command.Run(version, lambdas, false, role)

		assertExpectationsOnMocks(t)
	})
	t.Run("LambdaNameError", func(t *testing.T) {
		setupTest()

		finderErr := fmt.Errorf("error finding lambda names")
		finder.On("FindLambdas", lambdas).Return([]string{}, finderErr)

		printer.On("PrintErr", finderErr).Return()

		command.Run(version, lambdas, false, role)

		assertExpectationsOnMocks(t)
	})
}
