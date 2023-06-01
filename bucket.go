package promotions

import (
	"errors"
	"fmt"
	"os"
)

const bucket = "./bucket"

func (s Service) openFileFromBucket(name string) (*os.File, error) {
	path := fmt.Sprintf("%s/%s", bucket, name)
	if !fileExists(path) {
		return nil, errors.New("file not found")
	}
	return os.Open(path)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
