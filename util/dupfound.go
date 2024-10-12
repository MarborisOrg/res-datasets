package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Dataset struct {
	Tag       string   `json:"tag"`
	Patterns  []string `json:"patterns"`
	Responses []string `json:"responses"`
	Context   string   `json:"context"`
}

func loadDataset(filename string) ([]Dataset, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dataset []Dataset
	err = json.Unmarshal(file, &dataset)
	if err != nil {
		return nil, err
	}

	return dataset, nil
}

func findDuplicates(dataset []Dataset) map[string][]string {
	duplicatePatterns := make(map[string][]string)

	for _, data := range dataset {
		patternMap := make(map[string]bool)
		for _, pattern := range data.Patterns {
			if patternMap[pattern] {
				duplicatePatterns[data.Tag] = append(duplicatePatterns[data.Tag], pattern)
			} else {
				patternMap[pattern] = true
			}
		}
	}

	return duplicatePatterns
}

func logDuplicates(duplicates map[string][]string, logFile string) error {
	hasDuplicates := false
	for _, patterns := range duplicates {
		if len(patterns) > 0 {
			hasDuplicates = true
			break
		}
	}

	if !hasDuplicates {
		return nil
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer file.Close()

	logger := log.New(file, "duplicate-log: ", log.LstdFlags)
	for tag, patterns := range duplicates {
		if len(patterns) > 0 {
			logger.Printf("Tag: %s, Duplicates: %v\n", tag, patterns)
		}
	}

	return nil
}

func RunChecker() {
	dataset, err := loadDataset("./res/locales/en/intents.json")
	if err != nil {
		fmt.Printf("Error loading dataset: %v\n", err)
		return
	}

	duplicates := findDuplicates(dataset)

	for tag, patterns := range duplicates {
		if len(patterns) > 0 {
			fmt.Printf("Tag: %s, Duplicates: %v\n", tag, patterns)
		}
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("./log/duplicate_report_%s.log", currentTime)

	err = logDuplicates(duplicates, logFileName)
	if err != nil {
		fmt.Printf("Error logging duplicates: %v\n", err)
	} else if len(duplicates) == 0 {
		fmt.Println("No duplicates found, no log file created.")
	} else {
		fmt.Printf("Duplicate patterns logged to '%s'\n", logFileName)
	}
}
