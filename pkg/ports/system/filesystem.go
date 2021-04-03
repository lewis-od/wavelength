package system

import (
	"io/ioutil"
	"os"
)

type OSFilesystem struct {}

func (fs *OSFilesystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (fs *OSFilesystem) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
