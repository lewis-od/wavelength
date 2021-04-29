package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/io"
	"github.com/lewis-od/lambda-build/pkg/terraform"
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

type mockFilesystem struct {
	mock.Mock
}

func (m *mockFilesystem) ReadDir(dirname string) ([]io.FileInfo, error) {
	args := m.Called(dirname)
	return args.Get(0).([]io.FileInfo), args.Error(1)
}

func (m *mockFilesystem) FileExists(filename string) bool {
	args := m.Called(filename)
	return args.Bool(0)
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

func TestRun_Success(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	bucketName := "some-bucket"

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

	filesystem := new(mockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(version, lambdas, false)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_SkipBuild(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	bucketName := "some-bucket"

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

	filesystem := new(mockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(version, lambdas, true)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_OrchestratorError(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	bucketName := "bucket-name"

	orchestrator := new(mockOrchestrator)
	err := fmt.Errorf("Error text")
	orchestrator.On("BuildLambdas", lambdas).Return(err)

	tf := new(mockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	filesystem := new(mockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
	printer.On("Printlnf", "🪣 Found artifact bucket %s", []interface{}{bucketName}).Return()
	printer.On("PrintErr", err).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(version, lambdas, false)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_TerraformError(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}

	orchestrator := new(mockOrchestrator)

	tf := new(mockTerraform)
	err := fmt.Errorf("error")
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{}, err)

	filesystem := new(mockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(mockPrinter)
	printer.On("Printlnf", "🏗  Orchestrating upload of version %s of %s", []interface{}{version, lambdas}).Return()
	expectedError := fmt.Errorf("Could not determine name of artifact bucket from tf state\n%s", err)
	printer.On("PrintErr", expectedError).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(version, lambdas, false)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}
