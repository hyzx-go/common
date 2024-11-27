package utils

import (
	"encoding/json"
	"fmt"
)

// ToJSON converts a struct or map to a JSON string
func ToJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON parses a JSON string to the provided struct
func FromJSON(jsonStr string, v interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), v)
	if err != nil {
		return err
	}
	return nil
}

// PrettyPrintJSON prints JSON in a pretty format
func PrettyPrintJSON(v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Failed to generate JSON:", err)
		return
	}
	fmt.Println(string(bytes))
}
