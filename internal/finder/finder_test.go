package finder

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

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

var lambdasDir string = "lambdas"
var oneInfo io.FileInfo = io.FileInfo{Name: "one", IsDir: true}
var twoInfo io.FileInfo = io.FileInfo{Name: "two", IsDir: true}
var threeInfo io.FileInfo = io.FileInfo{Name: "three", IsDir: true}

func TestGetLambdas_NoArgs(t *testing.T) {
	filesystem := new(mockFilesystem)
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	finder := NewLambdaFinder(filesystem, lambdasDir)
	lambdas, err := finder.GetLambdas([]string{})

	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two"}, lambdas)
	filesystem.AssertExpectations(t)
}

func TestGetLambdas_ValidNames(t *testing.T) {
	filesystem := new(mockFilesystem)
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo, twoInfo, threeInfo}, nil)

	finder := NewLambdaFinder(filesystem, lambdasDir)
	lambdas, err := finder.GetLambdas([]string{"one", "two"})

	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two"}, lambdas)
	filesystem.AssertExpectations(t)
}

func TestGetLambdas_InvalidNameSupplied(t *testing.T) {
	filesystem := new(mockFilesystem)
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo}, nil)

	finder := NewLambdaFinder(filesystem, lambdasDir)
	lambdas, err := finder.GetLambdas([]string{"two"})

	assert.NotNil(t, err)
	assert.Nil(t, lambdas)
	filesystem.AssertExpectations(t)
}

func TestGetLambdas_FilesystemError(t *testing.T) {
	filesystem := new(mockFilesystem)
	fsError := fmt.Errorf("filesystem error")
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{}, fsError)

	finder := NewLambdaFinder(filesystem, lambdasDir)
	lambdas, err := finder.GetLambdas([]string{"one"})

	assert.Equal(t, fsError, err)
	assert.Nil(t, lambdas)
	filesystem.AssertExpectations(t)
}
