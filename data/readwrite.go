package data

import (
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func WriteToFile(filename string, content interface{}) error {
	data, err := yaml.Marshal(content)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func ReadFromFile(filename string, out interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var sb strings.Builder
	if _, err := io.Copy(&sb, file); err != nil {
		return err
	}

	if err := yaml.Unmarshal([]byte(sb.String()), out); err != nil {
		return err
	}

	return nil
}
