package helper

import (
	"encoding/json"
)

func Automapper(objOrigin interface{}, objDestination interface{}) {
	jsonOrigin := StructToJson(objOrigin)
	//log.Println("objOrigin:", objOrigin)
	json.Unmarshal([]byte(jsonOrigin), objDestination)
}
