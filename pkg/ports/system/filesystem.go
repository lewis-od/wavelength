package system

import (
	"github.com/lewis-od/lambda-build/pkg/io"
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
