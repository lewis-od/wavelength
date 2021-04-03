package builder

type Builder interface {
	BuildLambda(lambdaName string) error
	BuildLambdas(lambdaNames []string) error
}
