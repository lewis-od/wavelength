package builder

type BuildResult struct {
	LambdaName string
	Error      error
	Output     []byte
}

type Builder interface {
	BuildLambda(lambdaName string) *BuildResult
}
