package utils

import (
	log "github.com/sirupsen/logrus"
	"io/fs"
	"os"
)

type StringFilterFunction func(string) bool

var defaultStringFilterFunction StringFilterFunction = func(string) bool { return true }

// GetRecursiveFiles gets all files in a directory recursively
func GetRecursiveFiles(fsys fs.FS, root string, filterFunction StringFilterFunction) ([]string, error) {
	var files []string

	if filterFunction == nil {
		filterFunction = defaultStringFilterFunction
	}

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filterFunction(path) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// GetEnv gets an environment variable or panics if it doesn't exist
func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Error: environment variable %s not set", key)
	}
	return value
}
