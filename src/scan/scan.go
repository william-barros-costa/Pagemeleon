package scan

import (
	"errors"
	"fmt"
	"os"
)

func VerifyFile(fileLocation string) (bool, error) {
	fileInfo, err := os.Stat(fileLocation)
	if err != nil {
		return false, err
	}
	if fileInfo.Size() == 0 {
		return false, errors.New("File is empty")
	}
	if fileInfo.IsDir() {
		return false, fmt.Errorf("%q is a directory, not a file", fileLocation)
	}
	return true, nil
}
