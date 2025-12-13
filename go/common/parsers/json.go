package parsers

import (
	"encoding/json"
	"log"
	"os"
)

// Parse JSON files such as configs

func ParseJSONFile(filePath string, v any) (any, error) {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Could not read file %s: %v", filePath, err)
		return nil, err
	}

	// Unmarshal JSON content into the provided interface
	err = json.Unmarshal(content, v)
	if err != nil {
		log.Printf("Error unmarshalling JSON from file %s: %v", filePath, err)
		return nil, err
	}

	return v, nil
}
