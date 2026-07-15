package helper

import (
	"encoding/json"
)

func StructToJson(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func JsonToStruct(jsonStr string) (interface{}, error) {
	var data interface{}

	err := json.Unmarshal([]byte(jsonStr), &data)

	return data, err
}
