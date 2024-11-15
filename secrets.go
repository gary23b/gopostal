package gopostal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Secrets struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

func ReadSecrets() (*Secrets, error) {
	filePath, err := FindFileUpTheFolderTree("zzzSecrets.yaml", 4)
	if err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var ret Secrets

	err = yaml.Unmarshal(fileData, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func FindFileUpTheFolderTree(filename string, height int) (string, error) {
	// Make sure it is only a file name. Remove any directory info.
	_, filename = filepath.Split(filename)

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for i := 0; i <= height; i++ {
		filePath := filepath.Join(wd, filename)
		_, err := os.Stat(filePath)
		if errors.Is(err, os.ErrNotExist) {
			// Remove the last element of the working directory
			wd = filepath.Dir(wd)
			continue
		}
		return filePath, nil
	}
	return "", fmt.Errorf("file not found")
}
