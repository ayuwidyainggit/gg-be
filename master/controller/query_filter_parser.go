package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

// parseCSVIntValues parses comma-separated integers while preserving order and removing duplicates.
func parseCSVIntValues(rawValue string, fieldName string) ([]int, error) {
	rawValue = strings.TrimSpace(rawValue)
	if rawValue == "" {
		return nil, nil
	}
	if strings.EqualFold(rawValue, "null") || strings.EqualFold(rawValue, "undefined") {
		return nil, nil
	}

	parts := strings.Split(rawValue, ",")
	result := make([]int, 0, len(parts))
	seen := make(map[int]struct{})

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.EqualFold(part, "null") || strings.EqualFold(part, "undefined") {
			continue
		}

		value, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid %s", fieldName)
		}

		if _, exists := seen[value]; exists {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result, nil
}

// parseIntSliceQuery parses repeated and comma-separated integer query params.
func parseIntSliceQuery(args *fasthttp.Args, fieldName string, keys ...string) ([]int, error) {
	return parseIntSliceQueryWithOptions(args, fieldName, false, keys...)
}

// parseIntSliceQueryAllowZero parses repeated and comma-separated integer query params while preserving zero values.
func parseIntSliceQueryAllowZero(args *fasthttp.Args, fieldName string, keys ...string) ([]int, error) {
	return parseIntSliceQueryWithOptions(args, fieldName, true, keys...)
}

func parseIntSliceQueryWithOptions(args *fasthttp.Args, fieldName string, allowZero bool, keys ...string) ([]int, error) {
	result := make([]int, 0)
	seen := make(map[int]struct{})
	processedKeys := make(map[string]struct{})

	for _, key := range keys {
		processedKeys[key] = struct{}{}
		values := args.PeekMulti(key)
		for _, value := range values {
			parsedValues, err := parseCSVIntValues(string(value), fieldName)
			if err != nil {
				return nil, err
			}

			for _, parsedValue := range parsedValues {
				if fieldName == "distributor_id" && parsedValue <= 0 && !(allowZero && parsedValue == 0) {
					continue
				}

				if _, exists := seen[parsedValue]; exists {
					continue
				}

				seen[parsedValue] = struct{}{}
				result = append(result, parsedValue)
			}
		}
	}

	args.VisitAll(func(key, value []byte) {
		keyString := string(key)
		if _, exists := processedKeys[keyString]; exists {
			return
		}
		if !hasIndexedArrayPrefix(keyString, keys) {
			return
		}

		parsedValues, err := parseCSVIntValues(string(value), fieldName)
		if err != nil {
			result = nil
			seen = nil
			processedKeys = nil
			return
		}

		for _, parsedValue := range parsedValues {
			if fieldName == "distributor_id" && parsedValue <= 0 && !(allowZero && parsedValue == 0) {
				continue
			}

			if _, exists := seen[parsedValue]; exists {
				continue
			}

			seen[parsedValue] = struct{}{}
			result = append(result, parsedValue)
		}
	})

	if seen == nil {
		return nil, fmt.Errorf("invalid %s", fieldName)
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

func hasIndexedArrayPrefix(actualKey string, keys []string) bool {
	for _, key := range keys {
		trimmedKey := strings.TrimSuffix(key, "[]")
		if strings.HasPrefix(actualKey, trimmedKey+"[") && strings.HasSuffix(actualKey, "]") {
			return true
		}
	}

	return false
}

// parseStringSliceQuery parses repeated and comma-separated string query params.
func parseStringSliceQuery(args *fasthttp.Args, keys ...string) []string {
	result := make([]string, 0)
	seen := make(map[string]struct{})

	for _, key := range keys {
		values := args.PeekMulti(key)
		for _, value := range values {
			parts := strings.Split(string(value), ",")
			for _, part := range parts {
				cleaned := strings.TrimSpace(part)
				if cleaned == "" {
					continue
				}

				if _, exists := seen[cleaned]; exists {
					continue
				}

				seen[cleaned] = struct{}{}
				result = append(result, cleaned)
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func containsInt(values []int, expected int) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
