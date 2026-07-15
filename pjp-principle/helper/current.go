package helper

import "github.com/gin-gonic/gin"

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
