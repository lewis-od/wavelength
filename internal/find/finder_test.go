package find

import (
	"fmt"
	"github.com/lewis-od/lambda-build/internal/io"
	"github.com/lewis-od/lambda-build/internal/terraform"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_filesystem"
	"github.com/lewis-od/lambda-build/internal/testutil/mock_terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var lambdasDir = "lambdas"
var artifactStorageComponent = "terraform/artifact-storage"
var oneInfo = io.FileInfo{Name: "one", IsDir: true}
var twoInfo = io.FileInfo{Name: "two", IsDir: true}
var threeInfo = io.FileInfo{Name: "three", IsDir: true}
var bucketName = "my-bucket"

var filesystem *mock_filesystem.MockFilesystem
var tf *mock_terraform.MockTerraform

func initMocks() {
	filesystem = new(mock_filesystem.MockFilesystem)
	tf = new(mock_terraform.MockTerraform)
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
	).Return(map[string]terraform.Output{"bucket_name": {Value: bucketName}}, nil)

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
