package common

import (
	"encoding/json"
	"fmt"
)

type Serializable interface {
	ToJson() (string, error)
	ToMap() (map[string]any, error)
	GetID() (string, error)
	GetField(field string) (any, error)
}

// Default function for turning a serializable into JSON
func ToJson(object Serializable) (string, error) {
	jsonData, err := json.Marshal(object)
	jsonString := string(jsonData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return jsonString, err
	}

	return jsonString, nil
}

// Default function for turning serializable to Map
func ToMap(object Serializable) (map[string]any, error) {
	jsonData, err := object.ToJson()
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
