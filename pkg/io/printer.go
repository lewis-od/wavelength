package io

type Printer interface {
	Println(a ...interface{})
	Printlnf(format string, a ...interface{})
}
