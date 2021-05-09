package progress


type Action struct {
	InProgress string
	Success    string
	Error      string
}

var Build = Action{
	InProgress: "🔨 Building %s...\n",
	Success:    "✅ Building %s...done",
	Error:      "❌ Building %s...error",
}

var Upload = Action{
	InProgress: "☁️ Uploading %s...\n",
	Success:    "✅ Uploading %s...done",
	Error:      "❌ Uploading %s...error",
}

type BuildDisplay interface {
	Init(action Action)
	Started(lambdaName string)
	Completed(lambdaName string, wasSuccessful bool)
}
