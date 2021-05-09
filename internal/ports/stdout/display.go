package stdout

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/progress"
)

func NewBasicDisplay() progress.BuildDisplay {
	return &basicDisplay{
		action: progress.Build,
	}
}

type basicDisplay struct{
	action progress.Action
}

func (b *basicDisplay) Init(action progress.Action) {
	b.action = action
}

func (b *basicDisplay) Started(lambdaName string) {
	fmt.Printf(b.action.InProgress, lambdaName)
}

func (b *basicDisplay) Completed(lambdaName string, wasSuccessful bool) {
	if wasSuccessful {
		fmt.Printf(b.action.Success, lambdaName)
	} else {
		fmt.Printf(b.action.Error, lambdaName)
	}
}
