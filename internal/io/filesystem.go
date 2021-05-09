package io

type FileInfo struct {
	Name  string
	IsDir bool
}

type Filesystem interface {
	ReadDir(dirname string) ([]FileInfo, error)
	FileExists(filename string) bool
	AppendToFile(location string, contents string) error
}
