package system

import (
	"github.com/lewis-od/wavelength/internal/io"
	"io/ioutil"
	"os"
)

type osFilesystem struct{}

func NewFilesystem() io.Filesystem {
	return &osFilesystem{}
}

func (fs *osFilesystem) ReadDir(dirname string) ([]io.FileInfo, error) {
	results, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	info := make([]io.FileInfo, len(results))
	for i, result := range results {
		info[i] = io.FileInfo{
			Name:  result.Name(),
			IsDir: result.IsDir(),
		}
	}
	return info, nil
}

func (fs *osFilesystem) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

const userWritePermissions os.FileMode = 0644

func (fs *osFilesystem) AppendToFile(location string, contents string) error {
	f, err := os.OpenFile(location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, userWritePermissions)
	defer f.Close()
	if err != nil {
		return err
	}
	if _, err = f.WriteString(contents); err != nil {
		return err
	}
	return nil
}
