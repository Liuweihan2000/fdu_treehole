package utils

import "encoding/json"

func Marshal(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}

func UnMarshal(jsonString string, v interface{}) error {
	err := json.Unmarshal([]byte(jsonString), v)
	return err
}