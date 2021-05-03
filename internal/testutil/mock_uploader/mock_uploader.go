package mock_uploader

import "github.com/stretchr/testify/mock"

type MockUploader struct {
	mock.Mock
}

func (m *MockUploader) UploadLambda(version, bucketName, lambdaName, artifactLocation string) error {
	args := m.Called(version, bucketName, lambdaName, artifactLocation)
	return args.Error(0)
}
