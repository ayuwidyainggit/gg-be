package middleware

import (
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/exception"
	"scyllax-pjp/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JwtMiddleware(ctx *gin.Context) {
	var token string
	authorizationHeader := ctx.Request.Header.Get("Authorization")
	fields := strings.Fields(authorizationHeader)

	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		panic(exception.NewUnauthorizedError("empty token"))

	}

	config, _ := config.LoadConfig(".")
	payload, err := utils.DecodeToken(token, config.TokenSecret)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	claims, ok := payload.(jwt.MapClaims)
	if !ok {
		panic(exception.NewUnauthorizedError("invalid token"))
	}

	customerId, ok := claims["cust_id"].(string)
	if !ok {
		panic(exception.NewUnauthorizedError("invalid token"))
	}
	if parentCustomerId, ok := claims["parent_cust_id"].(string); ok {
		ctx.Set("parentCustomerId", parentCustomerId)
	}
	empId, okEmp := claims["emp_id"]
	if okEmp {
		ctx.Set("currentEmpId", empId)
	}

	// Extract user_id from claims
	var userId int64
	if userIdRaw, exists := claims["user_id"]; exists {
		switch v := userIdRaw.(type) {
		case float64:
			userId = int64(v)
		case int64:
			userId = v
		case int:
			userId = int64(v)
		}
		ctx.Set("currentUserId", userId)
	}

	// Pengecekan jika "emp_code" ada di dalam claims
	var empCode string
	if empCodeRaw, exists := claims["emp_code"]; exists {
		empCode, ok = empCodeRaw.(string)
		if !ok {
			panic(exception.NewUnauthorizedError("invalid emp_code format"))
		}
		ctx.Set("empCode", empCode)
	} else {
		// Jika emp_code tidak ada, bisa diabaikan atau tambahkan kondisi lain
		ctx.Set("empCode", "") // Default kosong jika tidak ada
	}

	ctx.Set("currentCustomerId", customerId)
	ctx.Next()
}
