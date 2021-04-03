package io

import (
	"os"
)

type Filesystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
	FileExists(filename string) bool
}
