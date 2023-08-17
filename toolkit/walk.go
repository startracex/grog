package toolkit

import (
	"os"
	"path/filepath"
)

type WalkFunc = func(path string, fi os.FileInfo) bool

// WalkExt Example: Get a list of files with HTML extensions
// WalkExt(directory,".html")
func WalkExt(dir string, ext string) []string {
	return Walk(dir, func(path string, fi os.FileInfo) bool {
		return !fi.IsDir() && filepath.Ext(path) == ext
	})
}

func WalkFiles(dir string) []string {
	return Walk(dir, func(path string, fi os.FileInfo) bool {
		return !fi.IsDir()
	})
}

// Walk When f returns true, append filepath.Clean(path) to the result list
func Walk(dir string, f func(path string, f os.FileInfo) bool) []string {
	var fileList []string
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if f(path, fi) {
			fileList = append(fileList, filepath.Clean(path))
		}
		return nil
	})
	if err != nil {
		return []string{}
	}
	return fileList
}
