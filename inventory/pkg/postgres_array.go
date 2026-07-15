package pkg

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Int64Array []int64

// Value implements driver.Valuer interface
func (a Int64Array) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}

	// Convert []int64 to PostgreSQL array format: {1,2,3}
	strValues := make([]string, len(a))
	for i, v := range a {
		strValues[i] = fmt.Sprintf("%d", v)
	}
	return "{" + strings.Join(strValues, ",") + "}", nil
}

func (a *Int64Array) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("cannot scan %T into Int64Array", value)
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	if str == "" {
		*a = Int64Array{}
		return nil
	}

	parts := strings.Split(str, ",")
	result := make(Int64Array, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		var val int64
		if _, err := fmt.Sscanf(part, "%d", &val); err != nil {
			return fmt.Errorf("cannot parse %q as int64: %w", part, err)
		}
		result = append(result, val)
	}

	*a = result
	return nil
}
