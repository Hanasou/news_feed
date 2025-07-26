package common

import (
	"encoding/json"
	"fmt"
)

type Serializable interface {
	ToJson() (string, error)
	GetId() (string, error)
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
