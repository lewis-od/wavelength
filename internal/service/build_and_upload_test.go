package service

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/terraform"
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

type mockTerraform struct {
	mock.Mock
}

func (m *mockTerraform) Output(directory string) (map[string]terraform.Output, error) {
	args := m.Called(directory)
	return args.Get(0).(map[string]terraform.Output), args.Error(1)
}

type mockFinder struct {
	mock.Mock
}

func (m *mockFinder) GetLambdas(lambdaNames []string) ([]string, error) {
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

func TestRun_Success(t *testing.T) {
	orchestrator := new(mockOrchestrator)
	orchestrator.On(
		"BuildLambdas",
		lambdas,
	).Return(nil)
	orchestrator.On(
		"UploadLambdas",
		version, bucketName, lambdas,
	).Return(nil)

	tf := new(mockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	finder := new(mockFinder)
	finder.On("GetLambdas", lambdas).Return(lambdas, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

	command := NewBuildAndUploadService(orchestrator, tf, finder, printer)
	command.Run(version, lambdas, false)

	mock.AssertExpectationsForObjects(t, orchestrator, tf, finder, printer)
}

func TestRun_SkipBuild(t *testing.T) {
	orchestrator := new(mockOrchestrator)
	orchestrator.On(
		"UploadLambdas",
		version, bucketName, lambdas,
	).Return(nil)

	tf := new(mockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	finder := new(mockFinder)
	finder.On("GetLambdas", lambdas).Return(lambdas, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

	command := NewBuildAndUploadService(orchestrator, tf, finder, printer)
	command.Run(version, lambdas, true)

	mock.AssertExpectationsForObjects(t, orchestrator, tf, finder, printer)
}

func TestRun_OrchestratorError(t *testing.T) {
	orchestrator := new(mockOrchestrator)
	err := fmt.Errorf("Error text")
	orchestrator.On("BuildLambdas", lambdas).Return(err)

	tf := new(mockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	finder := new(mockFinder)
	finder.On("GetLambdas", lambdas).Return(lambdas, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
	printer.On("Printlnf", "🪣 Found artifact bucket %s", []interface{}{bucketName}).Return()
	printer.On("PrintErr", err).Return()

	command := NewBuildAndUploadService(orchestrator, tf, finder, printer)
	command.Run(version, lambdas, false)

	mock.AssertExpectationsForObjects(t, orchestrator, tf, finder, printer)
}

func TestRun_TerraformError(t *testing.T) {
	orchestrator := new(mockOrchestrator)

	tf := new(mockTerraform)
	err := fmt.Errorf("error")
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{}, err)

	finder := new(mockFinder)
	finder.On("GetLambdas", lambdas).Return(lambdas, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
	expectedError := fmt.Errorf("Could not determine name of artifact bucket from tf state\n%s", err)
	printer.On("PrintErr", expectedError).Return()

	command := NewBuildAndUploadService(orchestrator, tf, finder, printer)
	command.Run(version, lambdas, false)

	mock.AssertExpectationsForObjects(t, orchestrator, tf, finder, printer)
}

func TestRun_FinderError(t *testing.T) {
	orchestrator := new(mockOrchestrator)
	tf := new(mockTerraform)

	finder := new(mockFinder)
	finderErr := fmt.Errorf("finder error")
	finder.On("GetLambdas", lambdas).Return([]string{}, finderErr)

	printer := new(mockPrinter)
	printer.On("PrintErr", finderErr).Return()

	command := NewBuildAndUploadService(orchestrator, tf, finder, printer)
	command.Run(version, lambdas, false)

	mock.AssertExpectationsForObjects(t, orchestrator, tf, finder, printer)
}
