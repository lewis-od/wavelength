package builder

type Builder interface {
	BuildLambda(lambdaName string) ([]byte, error)
}
