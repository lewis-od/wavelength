package progress

type BuildDisplay interface {
	Started(lambdaName string)
	Completed(lambdaName string, wasSuccessful bool)
}
