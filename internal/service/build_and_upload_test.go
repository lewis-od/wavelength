package service

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_finder"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_orchestrator"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_printer"
	"github.com/stretchr/testify/mock"
	"testing"
)

var version string = "version"
var lambdas []string = []string{"one", "two"}
var bucketName string = "some-bucket"

var orchestrator *mock_orchestrator.MockOrchestrator
var finder *mock_finder.MockFinder
var printer *mock_printer.MockPrinter

func initMocks() {
	orchestrator = new(mock_orchestrator.MockOrchestrator)
	finder = new(mock_finder.MockFinder)
	printer = new(mock_printer.MockPrinter)
}

func assertExpectationsOnMocks(t *testing.T) {
	mock.AssertExpectationsForObjects(t, orchestrator, finder, printer)
}

func TestRun_Success(t *testing.T) {
	initMocks()
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

	command := NewBuildAndUploadService(orchestrator, finder, printer)
	command.Run(version, lambdas, false)

	assertExpectationsOnMocks(t)
}

func TestRun_SkipBuild(t *testing.T) {
	orchestrator.On(
		"UploadLambdas",
		version, bucketName, lambdas,
	).Return(nil)

	finder.On("FindLambdas", lambdas).Return(lambdas, nil)
	finder.On("FindArtifactBucketName").Return(bucketName, nil)

	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

	command := NewBuildAndUploadService(orchestrator, finder, printer)
	command.Run(version, lambdas, true)

	assertExpectationsOnMocks(t)
}

func TestRun_OrchestratorError(t *testing.T) {
	initMocks()
	err := fmt.Errorf("error text")
	orchestrator.On("BuildLambdas", lambdas).Return(err)

	finder.On("FindLambdas", lambdas).Return(lambdas, nil)
	finder.On("FindArtifactBucketName").Return(bucketName, nil)

	printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
	printer.On("Printlnf", "🪣 Found artifact bucket %s", []interface{}{bucketName}).Return()
	printer.On("PrintErr", err).Return()

	command := NewBuildAndUploadService(orchestrator, finder, printer)
	command.Run(version, lambdas, false)

	assertExpectationsOnMocks(t)
}

func TestRun_FindArtifactBucketError(t *testing.T) {
	initMocks()

	finder.On("FindLambdas", lambdas).Return(lambdas, nil)
	finderErr := fmt.Errorf("error finding artifact bucket")
	finder.On("FindArtifactBucketName").Return(bucketName, finderErr)

	printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
	printer.On("PrintErr", finderErr).Return()

	command := NewBuildAndUploadService(orchestrator, finder, printer)
	command.Run(version, lambdas, false)

	assertExpectationsOnMocks(t)
}

func TestRun_LambdaNameError(t *testing.T) {
	initMocks()

	finderErr := fmt.Errorf("error finding lambda names")
	finder.On("FindLambdas", lambdas).Return([]string{}, finderErr)

	printer.On("PrintErr", finderErr).Return()

	command := NewBuildAndUploadService(orchestrator, finder, printer)
	command.Run(version, lambdas, false)

	assertExpectationsOnMocks(t)
}
