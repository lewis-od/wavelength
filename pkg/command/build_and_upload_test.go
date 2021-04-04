package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/io"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockOrchestrator struct {
	mock.Mock
}

func (m *MockOrchestrator) RunBuild(version, bucketName string, specifiedLambdas []string) error {
	args := m.Called(version, bucketName, specifiedLambdas)
	return args.Error(0)
}

type MockTerraform struct {
	mock.Mock
}

func (m *MockTerraform) Output(directory string) (map[string]terraform.Output, error) {
	args := m.Called(directory)
	return args.Get(0).(map[string]terraform.Output), args.Error(1)
}

type MockFilesystem struct {
	mock.Mock
}

func (m *MockFilesystem) ReadDir(dirname string) ([]io.FileInfo, error) {
	args := m.Called(dirname)
	return args.Get(0).([]io.FileInfo), args.Error(1)
}

func (m *MockFilesystem) FileExists(filename string) bool {
	args := m.Called(filename)
	return args.Bool(0)
}

type MockPrinter struct {
	mock.Mock
}

func (n *MockPrinter) Println(a ...interface{}) {
	n.Called(a)
}

func (n *MockPrinter) Printlnf(format string, a ...interface{}) {
	n.Called(format, a)
}

func TestRun_Success(t *testing.T) {
	arguments := []string{"version", "one", "two"}
	bucketName := "some-bucket"

	orchestrator := new(MockOrchestrator)
	orchestrator.On(
		"RunBuild",
		arguments[0], bucketName, arguments[1:],
	).Return(nil)

	tf := new(MockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	filesystem := new(MockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(MockPrinter)
	printer.On("Printlnf", mock.Anything, mock.Anything).Return()
	printer.On("Printlnf", mock.Anything, mock.Anything, mock.Anything).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(arguments)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_OrchestratorError(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	arguments := append([]string{version}, lambdas...)
	bucketName := "bucket-name"

	orchestrator := new(MockOrchestrator)
	err := fmt.Errorf("Error text")
	orchestrator.On("RunBuild", arguments[0], bucketName, arguments[1:]).Return(err)

	tf := new(MockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	filesystem := new(MockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(MockPrinter)
	printer.On("Printlnf", "🏗  Building version %s of %s", []interface{}{version, lambdas}).Return()
	printer.On("Printlnf", "🪣 Found artifact bucket %s", []interface{}{bucketName}).Return()
	printer.On("Println", []interface{}{"❌", err}).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(arguments)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_TerraformError(t *testing.T) {
	version := "version"
	lambdas := []string{"one", "two"}
	arguments := append([]string{version}, lambdas...)

	orchestrator := new(MockOrchestrator)

	tf := new(MockTerraform)
	err := fmt.Errorf("error")
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{}, err)

	filesystem := new(MockFilesystem)
	oneInfo := io.FileInfo{Name: "one", IsDir: true}
	twoInfo := io.FileInfo{Name: "two", IsDir: true}
	filesystem.On("ReadDir", "lambdas").Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	printer := new(MockPrinter)
	printer.On("Printlnf", "🏗  Building version %s of %s", []interface{}{version, lambdas}).Return()
	expectedError := fmt.Errorf("Could not determine name of artifact bucket from tf state\n%s", err)
	printer.On("Println", []interface{}{"❌", expectedError}).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, filesystem, printer)
	command.Run(arguments)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	filesystem.AssertExpectations(t)
	printer.AssertExpectations(t)
}
