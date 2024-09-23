package helper

import (
	"os"
)

func DoesFileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}
