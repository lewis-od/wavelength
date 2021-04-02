package command

import (
	"fmt"
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

type MockLerna struct {
	mock.Mock
}

func (m *MockLerna) BuildLambda(lambdaName string) error {
	args := m.Called(lambdaName)
	return args.Error(0)
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

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", "lambda-one").Return(nil)
	mockLerna.On("BuildLambda", "lambda-two").Return(nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockFilesystem)
	cmd.Run([]string{})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-one")
	mockLerna.AssertCalled(t, "BuildLambda", "lambda-two")
	mockLerna.AssertNotCalled(t, "BuildLambda", "some-file")
}

func TestRun_AllLambdas_OneError(t *testing.T) {
	lambdaOneDirectory := &FakeFile{name: "lambda-one", isDir: true}
	lambdaTwoDirectory := &FakeFile{name: "lambda-two", isDir: true}
	otherFile := &FakeFile{name: "some-file", isDir: false}

	mockFilesystem := new(MockFilesystem)
	mockFilesystem.On(
		"ReadDir",
		"lambdas",
	).Return([]os.FileInfo{lambdaOneDirectory, lambdaTwoDirectory, otherFile}, nil)

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", "lambda-one").Return(fmt.Errorf("error"))
	mockLerna.On("BuildLambda", "lambda-two").Return(nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockFilesystem)
	cmd.Run([]string{})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-one")
	mockLerna.AssertNotCalled(t, "BuildLambda", "lambda-two")
	mockLerna.AssertNotCalled(t, "BuildLambda", "some-file")
}

func TestRun_SingleLambdas_Success(t *testing.T) {
	lambdaName := "lambda-name"
	lambdaOneDirectory := &FakeFile{name: lambdaName, isDir: true}

	mockFilesystem := new(MockFilesystem)
	mockFilesystem.On(
		"ReadDir",
		"lambdas",
	).Return([]os.FileInfo{lambdaOneDirectory}, nil)

	mockLerna := new(MockLerna)
	mockLerna.On("BuildLambda", lambdaName).Return(nil)

	cmd := NewBuildAndUploadCommand(mockLerna, mockFilesystem)
	cmd.Run([]string{lambdaName})

	mockLerna.AssertCalled(t, "BuildLambda", "lambda-name")
}
