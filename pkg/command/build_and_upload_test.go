package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

type MockFilesystem struct {
	mock.Mock
}

func (m *MockFilesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	args := m.Called(dirname)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

func (m *MockFilesystem) FileExists(filename string) bool {
	args := m.Called(filename)
	return args.Bool(0)
}

type MockTerraform struct {
	mock.Mock
}

func (m *MockTerraform) Output(directory string) (map[string]terraform.Output, error) {
	args := m.Called(directory)
	return args.Get(0).(map[string]terraform.Output), args.Error(1)
}

type MockOrchestrator struct {
	mock.Mock
}

func (m *MockOrchestrator) RunBuild(specifiedLambdas []string) error {
	args := m.Called(specifiedLambdas)
	return args.Error(0)
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
	arguments := []string{"one", "two"}

	orchestrator := new(MockOrchestrator)
	orchestrator.On("RunBuild", arguments).Return(nil)

	tf := new(MockTerraform)
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: "some-bucket"}}, nil)

	printer := new(MockPrinter)
	printer.On("Println", mock.Anything).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, printer)
	command.Run(arguments)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_OrchestratorError(t *testing.T) {
	arguments := []string{"one", "two"}

	orchestrator := new(MockOrchestrator)
	err := fmt.Errorf("Error text")
	orchestrator.On("RunBuild", arguments).Return(err)

	tf := new(MockTerraform)

	printer := new(MockPrinter)
	printer.On("Println", []interface{}{"❌", err}).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, printer)
	command.Run(arguments)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	printer.AssertExpectations(t)
}

func TestRun_TerraformError(t *testing.T) {
	arguments := []string{"one", "two"}

	orchestrator := new(MockOrchestrator)
	orchestrator.On("RunBuild", arguments).Return(nil)

	tf := new(MockTerraform)
	err := fmt.Errorf("error")
	tf.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{}, err)

	printer := new(MockPrinter)
	expectedError := fmt.Errorf("Could not determine name of artifact bucket from tf state\n%s", err)
	printer.On("Println", []interface{}{expectedError}).Return()

	command := NewBuildAndUploadCommand(orchestrator, tf, printer)
	command.Run(arguments)

	orchestrator.AssertExpectations(t)
	tf.AssertExpectations(t)
	printer.AssertExpectations(t)
}
