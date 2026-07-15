package helper

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCurrentCustomerId(ctx *gin.Context) (string, bool) {
	currentCustomerId, exists := ctx.Get("currentCustomerId")
	if !exists {
		return "", false
	}
	customerId, ok := currentCustomerId.(string)
	if !ok {
		return "", false
	}
	return customerId, true
}

func GetParentCustomerId(ctx *gin.Context) (string, bool) {
	parentCustomerId, exists := ctx.Get("parentCustomerId")
	if !exists {
		return "", false
	}
	customerId, ok := parentCustomerId.(string)
	if !ok || customerId == "" {
		return "", false
	}
	return customerId, true
}

func GetCurrentUserId(ctx *gin.Context) (int64, bool) {
	currentUserId, exists := ctx.Get("currentUserId")
	if !exists {
		return 0, false
	}
	userId, ok := currentUserId.(int64)
	if !ok {
		return 0, false
	}
	return userId, true
}

func GetCurrentEmpId(ctx *gin.Context) (int64, bool) {
	currentEmpId, exists := ctx.Get("currentEmpId")
	if !exists {
		return 0, false
	}

	switch v := currentEmpId.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, false
		}
		return n, true
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}
