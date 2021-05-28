package builder

type Updater interface {
	UpdateCode(lambdaName, bucketName, bucketKey string, role *AssumeRole) error
}
