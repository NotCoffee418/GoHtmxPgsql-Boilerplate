package utils

import (
	"os"
	"path/filepath"
)

type StringFilterFunction func(string) bool

var defaultStringFilterFunction StringFilterFunction = func(string) bool { return true }

func GetRecursiveFiles(dir string, filterFunction StringFilterFunction) ([]string, error) {
	var files []string

	if filterFunction == nil {
		filterFunction = defaultStringFilterFunction
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filterFunction(path) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
