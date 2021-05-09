package stdout

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/progress"
)

func NewBasicDisplay() progress.BuildDisplay {
	return &basicDisplay{}
}

type basicDisplay struct{}

func (b *basicDisplay) Started(lambdaName string) {
	fmt.Printf("🔨 Building %s...\n", lambdaName)
}

func (b *basicDisplay) Completed(lambdaName string, wasSuccessful bool) {
	if wasSuccessful {
		fmt.Printf("✅ Successfully built %s\n", lambdaName)
	} else {
		fmt.Printf("❌ Error building %s\n", lambdaName)
	}
}
