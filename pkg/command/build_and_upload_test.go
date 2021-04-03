package command

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/terraform"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
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

type MockLerna struct {
	mock.Mock
}

func (m *MockLerna) BuildLambda(lambdaName string) error {
	args := m.Called(lambdaName)
	return args.Error(0)
}

type MockTerraform struct {
	mock.Mock
}

func (m *MockTerraform) Output(directory string) (map[string]terraform.Output, error) {
	args := m.Called(directory)
	return args.Get(0).(map[string]terraform.Output), args.Error(1)
}

type FakeFile struct {
	name string
	isDir bool
}

func (f *FakeFile) Name() string {
	return f.name
}

func (f *FakeFile) Size() int64 {
	return 1
}

func (f *FakeFile) Mode() os.FileMode {
	return os.ModeDir
}

func (f *FakeFile) ModTime() time.Time {
	return time.Now()
}

func (f *FakeFile) IsDir() bool {
	return f.isDir
}

func (f *FakeFile) Sys() interface{} {
	return ""
}

func TestRun_AllLambdas_Success(t *testing.T) {
	lambdaOneDirectory := &FakeFile{name: "lambda-one", isDir: true}
	lambdaTwoDirectory := &FakeFile{name: "lambda-two", isDir: true}
	otherFile := &FakeFile{name: "some-file", isDir: false}

	mockFilesystem := new(MockFilesystem)
	mockFilesystem.On(
		"ReadDir",
		"lambdas",
	).Return([]os.FileInfo{lambdaOneDirectory, lambdaTwoDirectory, otherFile}, nil)
	mockFilesystem.On(
		"FileExists",
		"lambdas/lambda-one/dist/lambda-one.zip",
	).Return(true)
	mockFilesystem.On(
		"FileExists",
		"lambdas/lambda-two/dist/lambda-two.zip",
	).Return(true)

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", "lambda-one").Return(nil)
	mockLerna.On("BuildLambda", "lambda-two").Return(nil)

	mockTerraform := &MockTerraform{}
	mockTerraform.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{
		"bucket_name": terraform.Output{Value: "bucket"},
	}, nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockTerraform, mockFilesystem)
	cmd.Run([]string{})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-one")
	mockLerna.AssertCalled(t, "BuildLambda", "lambda-two")
	mockLerna.AssertNotCalled(t, "BuildLambda", "some-file")
}

func TestRun_AllLambdas_BuildError(t *testing.T) {
	lambdaOneDirectory := &FakeFile{name: "lambda-one", isDir: true}
	lambdaTwoDirectory := &FakeFile{name: "lambda-two", isDir: true}
	otherFile := &FakeFile{name: "some-file", isDir: false}

	mockFilesystem := new(MockFilesystem)
	mockFilesystem.On(
		"ReadDir",
		"lambdas",
	).Return([]os.FileInfo{lambdaOneDirectory, lambdaTwoDirectory, otherFile}, nil)
	mockFilesystem.On(
		"FileExists",
		"lambdas/lambda-one/dist/lambda-one.zip",
	).Return(true)

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", "lambda-one").Return(fmt.Errorf("error"))
	mockLerna.On("BuildLambda", "lambda-two").Return(nil)

	mockTerraform := &MockTerraform{}
	mockTerraform.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{
		"bucket_name": terraform.Output{Value: "bucket"},
	}, nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockTerraform, mockFilesystem)
	cmd.Run([]string{})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-one")
	mockLerna.AssertNotCalled(t, "BuildLambda", "lambda-two")
	mockLerna.AssertNotCalled(t, "BuildLambda", "some-file")
}

func TestRun_AllLambdas_ArtifactNotFound(t *testing.T) {
	lambdaOneDirectory := &FakeFile{name: "lambda-one", isDir: true}
	lambdaTwoDirectory := &FakeFile{name: "lambda-two", isDir: true}

	mockFilesystem := new(MockFilesystem)
	mockFilesystem.On(
		"ReadDir",
		"lambdas",
	).Return([]os.FileInfo{lambdaOneDirectory, lambdaTwoDirectory}, nil)
	mockFilesystem.On(
		"FileExists",
		"lambdas/lambda-one/dist/lambda-one.zip",
	).Return(false)

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", "lambda-one").Return(nil)

	mockTerraform := &MockTerraform{}
	mockTerraform.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{
		"bucket_name": terraform.Output{Value: "bucket"},
	}, nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockTerraform, mockFilesystem)
	cmd.Run([]string{})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-one")
	mockLerna.AssertNotCalled(t, "BuildLambda", "lambda-two")
}

func TestRun_SingleLambda_Success(t *testing.T) {
	lambdaName := "lambda-name"
	lambdaOneDirectory := &FakeFile{name: lambdaName, isDir: true}

	mockFilesystem := new(MockFilesystem)
	mockFilesystem.On(
		"ReadDir",
		"lambdas",
	).Return([]os.FileInfo{lambdaOneDirectory}, nil)
	mockFilesystem.On(
		"FileExists",
		"lambdas/lambda-name/dist/lambda-name.zip",
	).Return(true)

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", lambdaName).Return(nil)

	mockTerraform := &MockTerraform{}
	mockTerraform.On(
		"Output",
		"terraform/deployments/artifact-storage",
	).Return(map[string]terraform.Output{
		"bucket_name": terraform.Output{Value: "bucket"},
	}, nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockTerraform, mockFilesystem)
	cmd.Run([]string{lambdaName})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-name")
}
