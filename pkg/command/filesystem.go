package command

import (
	"io/ioutil"
	"os"
)

type Filesystem interface {
	ReadDir(dirname string) ([]os.FileInfo, error)
}

type OSFilesystem struct {}

func (fs *OSFilesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

