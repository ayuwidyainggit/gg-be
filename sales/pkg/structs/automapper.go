package structs

import (
	"encoding/json"
)

func Automapper(objOrigin interface{}, objDestination interface{}) error {
	var (
		err error
	)
	jsonOrigin := StructToJson(objOrigin)
	err = json.Unmarshal([]byte(jsonOrigin), &objDestination)
	return err
}
