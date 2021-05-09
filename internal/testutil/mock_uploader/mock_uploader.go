package mock_uploader

import (
	"github.com/lewis-od/wavelength/internal/builder"
	"github.com/stretchr/testify/mock"
)

type MockUploader struct {
	mock.Mock
}

func (m *MockUploader) UploadLambda(version, bucketName, lambdaName, artifactLocation string) *builder.BuildResult {
	args := m.Called(version, bucketName, lambdaName, artifactLocation)
	return args.Get(0).(*builder.BuildResult)
}
