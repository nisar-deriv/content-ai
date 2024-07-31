package data

import (
	"io"
	"os"
	"strings"
)

func WriteToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(content + "\n"); err != nil {
		return err
	}
	return nil
}

func ReadFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	if _, err := io.Copy(&sb, file); err != nil {
		return "", err
	}

	return sb.String(), nil
}
