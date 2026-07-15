package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/exception"
	"scyllax-pjp/utils"
	"strings"
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
