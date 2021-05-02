package service

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockOrchestrator struct {
	mock.Mock
}

func (m *mockOrchestrator) BuildLambdas(lambdas []string) error {
	args := m.Called(lambdas)
	return args.Error(0)
}

func (m *mockOrchestrator) UploadLambdas(version, bucketName string, lambdas []string) error {
	args := m.Called(version, bucketName, lambdas)
	return args.Error(0)
}

type mockFinder struct {
	mock.Mock
}

func (m *mockFinder) FindArtifactBucketName() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockFinder) FindLambdas(lambdaNames []string) ([]string, error) {
	args := m.Called(lambdaNames)
	return args.Get(0).([]string), args.Error(1)
}

type mockPrinter struct {
	mock.Mock
}

func (n *mockPrinter) Println(a ...interface{}) {
	n.Called(a)
}

func (n *mockPrinter) Printlnf(format string, a ...interface{}) {
	n.Called(format, a)
}

func (n *mockPrinter) PrintErr(err error) {
	n.Called(err)
}

var version string = "version"
var lambdas []string = []string{"one", "two"}
var bucketName string = "some-bucket"

var orchestrator *mockOrchestrator
var finder *mockFinder
var printer *mockPrinter

func initMocks() {
	orchestrator = new(mockOrchestrator)
	finder = new(mockFinder)
	printer = new(mockPrinter)
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
