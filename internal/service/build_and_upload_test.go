package service

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_finder"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_orchestrator"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_printer"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestBuildAndUploadService_Run(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	bucketName := "some-bucket"

	var orchestrator *mock_orchestrator.MockOrchestrator
	var finder *mock_finder.MockFinder
	var printer *mock_printer.MockPrinter
	var command BuildAndUploadService

	setupTest := func() {
		orchestrator = new(mock_orchestrator.MockOrchestrator)
		finder = new(mock_finder.MockFinder)
		printer = new(mock_printer.MockPrinter)
		command = NewBuildAndUploadService(orchestrator, finder, printer)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, orchestrator, finder, printer)
	}

	t.Run("Success", func(t *testing.T) {
		setupTest()
		orchestrator.On(
			"BuildLambdas",
			lambdas,
		).Return(nil)
		orchestrator.On(
			"UploadLambdas",
			version, bucketName, lambdas,
		).Return(nil)

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On("Printlnf", mock.Anything, mock.Anything).Return()
		printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

		command.Run(version, lambdas, false)

		assertExpectationsOnMocks(t)
	})
	t.Run("SkipBuild", func(t *testing.T) {
		setupTest()
		orchestrator.On(
			"UploadLambdas",
			version, bucketName, lambdas,
		).Return(nil)

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On("Printlnf", mock.Anything, mock.Anything).Return()
		printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

		command.Run(version, lambdas, true)

		assertExpectationsOnMocks(t)
	})
	t.Run("OrchestrationError", func(t *testing.T) {
		setupTest()
		err := fmt.Errorf("error text")
		orchestrator.On("BuildLambdas", lambdas).Return(err)

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finder.On("FindArtifactBucketName").Return(bucketName, nil)

		printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
		printer.On("Printlnf", "🪣 Found artifact bucket %s", []interface{}{bucketName}).Return()
		printer.On("PrintErr", err).Return()

		command.Run(version, lambdas, false)

		assertExpectationsOnMocks(t)
	})
	t.Run("FindArtifactBucketError", func(t *testing.T) {
		setupTest()

		finder.On("FindLambdas", lambdas).Return(lambdas, nil)
		finderErr := fmt.Errorf("error finding artifact bucket")
		finder.On("FindArtifactBucketName").Return(bucketName, finderErr)

		printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
		printer.On("PrintErr", finderErr).Return()

		command.Run(version, lambdas, false)

		assertExpectationsOnMocks(t)
	})
	t.Run("LambdaNameError", func(t *testing.T) {
		setupTest()

		finderErr := fmt.Errorf("error finding lambda names")
		finder.On("FindLambdas", lambdas).Return([]string{}, finderErr)

		printer.On("PrintErr", finderErr).Return()

		command.Run(version, lambdas, false)

		assertExpectationsOnMocks(t)
	})
}
