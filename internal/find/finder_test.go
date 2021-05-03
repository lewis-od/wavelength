package find

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/io"
	"github.com/lewis-od/wavelength/internal/terraform"
	"github.com/lewis-od/wavelength/internal/testutil/mock_filesystem"
	"github.com/lewis-od/wavelength/internal/testutil/mock_terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLambdaFinder(t *testing.T) {
	lambdasDir := "lambdas"
	artifactStorageComponent := "terraform/artifact-storage"
	outputName := "bucket_name"

	var filesystem *mock_filesystem.MockFilesystem
	var tf *mock_terraform.MockTerraform
	var finder Finder

	setupTest := func() {
		filesystem = new(mock_filesystem.MockFilesystem)
		tf = new(mock_terraform.MockTerraform)
		finder = NewLambdaFinder(filesystem, tf, &lambdasDir, &artifactStorageComponent, &outputName)
	}

	assertExpectationsOnMocks := func(t *testing.T) {
		mock.AssertExpectationsForObjects(t, filesystem, tf)
	}

	t.Run("GetLambdas", func(t *testing.T) {
		oneInfo := io.FileInfo{Name: "one", IsDir: true}
		twoInfo := io.FileInfo{Name: "two", IsDir: true}
		threeInfo := io.FileInfo{Name: "three", IsDir: true}

		t.Run("NoArgs", func(t *testing.T) {
			setupTest()
			filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo, twoInfo}, nil)

			lambdas, err := finder.FindLambdas([]string{})

			assert.Nil(t, err)
			assert.Equal(t, []string{"one", "two"}, lambdas)
			assertExpectationsOnMocks(t)
		})
		t.Run("ValidNames", func(t *testing.T) {
			setupTest()
			filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo, twoInfo, threeInfo}, nil)

			lambdas, err := finder.FindLambdas([]string{"one", "two"})

			assert.Nil(t, err)
			assert.Equal(t, []string{"one", "two"}, lambdas)
			assertExpectationsOnMocks(t)
		})
		t.Run("InvalidNameSupplied", func(t *testing.T) {
			setupTest()
			filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{oneInfo}, nil)

			lambdas, err := finder.FindLambdas([]string{"two"})

			assert.NotNil(t, err)
			assert.Nil(t, lambdas)
			assertExpectationsOnMocks(t)
		})
		t.Run("FilesystemError", func(t *testing.T) {
			setupTest()
			fsError := fmt.Errorf("filesystem error")
			filesystem.On("ReadDir", lambdasDir).Return([]io.FileInfo{}, fsError)

			lambdas, err := finder.FindLambdas([]string{"one"})

			assert.Equal(t, fsError, err)
			assert.Nil(t, lambdas)
			assertExpectationsOnMocks(t)
		})
	})
	t.Run("FindArtifactBucketName", func(t *testing.T) {
		bucketName := "my-bucket"

		t.Run("Success", func(t *testing.T) {
			setupTest()
			tf.On(
				"Output",
				artifactStorageComponent,
			).Return(map[string]terraform.Output{outputName: {Value: bucketName}}, nil)

			foundName, err := finder.FindArtifactBucketName()

			assert.Nil(t, err)
			assert.Equal(t, bucketName, foundName)
			assertExpectationsOnMocks(t)
		})
		t.Run("TerraformError", func(t *testing.T) {
			setupTest()
			tfError := fmt.Errorf("terraform error")
			tf.On(
				"Output",
				artifactStorageComponent,
			).Return(map[string]terraform.Output{}, tfError)

			foundName, err := finder.FindArtifactBucketName()

			assert.Equal(t, tfError, err)
			assert.Empty(t, foundName)
			assertExpectationsOnMocks(t)
		})
	})
}
