package handlers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Update structure to hold the weekly update information
type Update struct {
	Team     string   `yaml:"Team"`
	Problems []string `yaml:"Problems"`
	Progress []string `yaml:"Progress"`
	Insights []string `yaml:"Insights"`
	Plans    []string `yaml:"Plans"`
}

// Function to read and parse the content of the team file
func parseTeamFile(filePath string) (Update, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Update{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var update Update
	var currentSection *[]string

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "Team:"):
			update.Team = strings.TrimSpace(strings.TrimPrefix(line, "Team:"))
		case strings.HasPrefix(line, "Problems"):
			currentSection = &update.Problems
		case strings.HasPrefix(line, "Progress"):
			currentSection = &update.Progress
		case strings.HasPrefix(line, "Insights"):
			currentSection = &update.Insights
		case strings.HasPrefix(line, "Plans"):
			currentSection = &update.Plans
		case strings.HasPrefix(line, "•"):
			if currentSection != nil {
				*currentSection = append(*currentSection, strings.TrimSpace(strings.TrimPrefix(line, "•")))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Update{}, err
	}

	return update, nil
}

// Function to convert the update structure to YAML and write it to a file
func writeYAML(update Update, filePath string) error {
	data, err := yaml.Marshal(&update)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Main function to read all team files and convert them to YAML
func ConvertFilesToYaml() {
	weekFolder := getWeekFolder()
	files, err := os.ReadDir(weekFolder)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".txt") {
			filePath := filepath.Join(weekFolder, file.Name())
			update, err := parseTeamFile(filePath)
			if err != nil {
				fmt.Println("Error parsing file:", err)
				continue
			}

			yamlFilePath := filepath.Join(weekFolder, update.Team+".yaml")
			err = writeYAML(update, yamlFilePath)
			if err != nil {
				fmt.Println("Error writing YAML file:", err)
			}
			// Delete the original .txt file after successful conversion
			err = os.Remove(filePath)
			if err != nil {
				fmt.Println("Error deleting original file:", err)
			}
		}
	}

	fmt.Println("Conversion completed.")
}
