package utils

import (
	"os"
	"path/filepath"
)

/*
FindProjectRoot 查找项目的根目录

返回值

	string: 项目的根目录路径
*/
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		_, err := os.Stat(
			filepath.Join(dir, "go.mod"),
		)

		if err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}
