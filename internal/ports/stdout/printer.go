package stdout

import (
	"fmt"
	"github.com/lewis-od/wavelength/internal/io"
)

type printer struct{}

func NewPrinter() io.Printer {
	return &printer{}
}

func (p *printer) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (p *printer) Printlnf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func (p *printer) PrintErr(err error) {
	p.Println("❌", err)
}
