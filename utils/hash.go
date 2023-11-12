package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func HashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", nil
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func GetTargetFilePath(filePath string) (string, string, error) {
	hash, err := HashFile(filePath)
	if err != nil {
		return "", "", err
	}
	return hash[:2], hash[2:] + strings.ToLower(path.Ext(filePath)), nil
}

// PrepareTargetFolder creates the target folder for a file
// and returns the relative path to the file in the target folder.
func PrepareTargetFolder(filePath string, rootDir string) (string, error) {

	folder, fileName, err := GetTargetFilePath(filePath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path.Join(rootDir, folder), 0755)
	if err != nil {
		return "", err
	}

	return path.Join(folder, fileName), nil
}
