package builder

type Builder interface {
	BuildLambda(lambdaName string) error
}
