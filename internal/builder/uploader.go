package builder

type Uploader interface {
	UploadLambda(version, bucketName, lambdaName, artifactLocation string, role *Role) *BuildResult
}
