package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

// LookUpFilePath look up find file path by file name,
// maxDeep represents the maximum depth of the recursion
func LookUpFilePath(destName string, maxDeep int) (string, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var find func(dir string, deep int) string
	find = func(dir string, deep int) string {
		if deep > maxDeep {
			return ""
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return ""
		}

		parent := path.Dir(dir)
		for _, file := range files {
			if file.Mode()&os.ModeSymlink != 0 {
				continue
			}

			if !file.IsDir() {
				if file.Name() == destName {
					return dir
				}
			}

		}
		destDir := find(parent, deep+1)
		if destDir != "" {
			return destDir
		}

		return ""
	}

	dir := find(currentDir, 0)
	if dir == "" {
		return "", errors.New("not found")
	}
	return dir + "/", nil
}
