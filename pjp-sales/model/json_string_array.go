package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONStringArray []string

func (j JSONStringArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return []byte("[]"), nil
	}

	data, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (j *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*j = JSONStringArray{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported JSONStringArray scan type %T", value)
	}

	if len(bytes) == 0 {
		*j = JSONStringArray{}
		return nil
	}

	var result []string
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = JSONStringArray(result)
	return nil
}
