package find

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/lewis-od/lambda-build/internal/terraform"
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

type mockTerraform struct {
	mock.Mock
}

func (m *mockTerraform) Output(directory string) (map[string]terraform.Output, error) {
	args := m.Called(directory)
	return args.Get(0).(map[string]terraform.Output), args.Error(1)
}

var lambdasDir string = "lambdas"
var artifactStorageComponent = "terraform/artifact-storage"
var oneInfo io.FileInfo = io.FileInfo{Name: "one", IsDir: true}
var twoInfo io.FileInfo = io.FileInfo{Name: "two", IsDir: true}
var threeInfo io.FileInfo = io.FileInfo{Name: "three", IsDir: true}
var bucketName string = "my-bucket"

var filesystem *mockFilesystem
var tf *mockTerraform

func initMocks() {
	filesystem = new(mockFilesystem)
	tf = new(mockTerraform)
}

func assertExpectationsOnMocks(t *testing.T) {
	mock.AssertExpectationsForObjects(t, filesystem, tf)
}

func TestGetLambdas_NoArgs(t *testing.T) {
	initMocks()
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo, twoInfo}, nil)

	finder := NewLambdaFinder(filesystem, tf, lambdasDir, artifactStorageComponent)
	lambdas, err := finder.FindLambdas([]string{})

	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two"}, lambdas)
	assertExpectationsOnMocks(t)
}

func TestGetLambdas_ValidNames(t *testing.T) {
	initMocks()
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo, twoInfo, threeInfo}, nil)

	finder := NewLambdaFinder(filesystem, tf, lambdasDir, artifactStorageComponent)
	lambdas, err := finder.FindLambdas([]string{"one", "two"})

	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two"}, lambdas)
	assertExpectationsOnMocks(t)
}

func TestGetLambdas_InvalidNameSupplied(t *testing.T) {
	initMocks()
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo}, nil)

	finder := NewLambdaFinder(filesystem, tf, lambdasDir, artifactStorageComponent)
	lambdas, err := finder.FindLambdas([]string{"two"})

	assert.NotNil(t, err)
	assert.Nil(t, lambdas)
	assertExpectationsOnMocks(t)
}

func TestGetLambdas_FilesystemError(t *testing.T) {
	initMocks()
	fsError := fmt.Errorf("filesystem error")
	filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{}, fsError)

	finder := NewLambdaFinder(filesystem, tf, lambdasDir, artifactStorageComponent)
	lambdas, err := finder.FindLambdas([]string{"one"})

	assert.Equal(t, fsError, err)
	assert.Nil(t, lambdas)
	assertExpectationsOnMocks(t)
}

func TestFindArtifactBucketName_Success(t *testing.T) {
	initMocks()
	tf.On(
		"Output",
		artifactStorageComponent,
	).Return(map[string]terraform.Output{"bucket_name": terraform.Output{Value: bucketName}}, nil)

	finder := NewLambdaFinder(filesystem, tf, lambdasDir, artifactStorageComponent)
	foundName, err := finder.FindArtifactBucketName()

	assert.Nil(t, err)
	assert.Equal(t, bucketName, foundName)
	assertExpectationsOnMocks(t)
}

func TestFindArtifactBucketName_TerraformError(t *testing.T) {
	initMocks()
	tfError := fmt.Errorf("terraform error")
	tf.On(
		"Output",
		artifactStorageComponent,
	).Return(map[string]terraform.Output{}, tfError)

	finder := NewLambdaFinder(filesystem, tf, lambdasDir, artifactStorageComponent)
	foundName, err := finder.FindArtifactBucketName()

	assert.Equal(t, tfError, err)
	assert.Empty(t, foundName)
	assertExpectationsOnMocks(t)
}
