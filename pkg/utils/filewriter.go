package utils

import (
	"os"
	// "path/filepath"
)

func WriteLLMResponseToTempFile(data []byte, prefix string) (string, error) {
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, prefix+"*.json")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write(data)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}