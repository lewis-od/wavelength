package stdout

import (
	"fmt"
	"github.com/lewis-od/lambda-build/pkg/io"
)

type printer struct {}

func NewPrinter() io.Printer {
	return &printer{}
}

func (p *printer) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (p *printer) Printlnf(format string, a ...interface{}) {
	fmt.Printf(format + "\n", a...)
}
