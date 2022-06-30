package s3yandex

import (
	"os"
	"strings"
)

func GetFileBytes(path string, filename string) ([]byte, error) {
	if _, err := os.Stat(path + filename); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	fileBytes, err := os.ReadFile(path + filename)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}
	return fileBytes, nil
}
func GetFileNameByExt(file string, ext string) (string, bool) {
	split := strings.Split(file, ext)
	splitBySlash := strings.Split(strings.Join(split, ""), "/")
	afterSlash := splitBySlash[len(splitBySlash)-1]
	if strings.TrimSpace(afterSlash) == "" {
		return "", false
	}
	if l := len(strings.Split(afterSlash, ".")); l > 1 {
		return "", false
	}
	return afterSlash + ext, true
}
