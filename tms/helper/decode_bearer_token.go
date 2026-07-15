package helper

import (
	"errors"
	"scyllax-tms/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims sesuai dengan struktur token kamu
type CustomClaims struct {
	CustID string `json:"cust_id"`
	jwt.RegisteredClaims
}

// DecodeBearerToken menerima string Authorization dan secret, lalu mengembalikan cust_id
func DecodeBearerToken(ctx *fiber.Ctx) (string, error) {
	authHeader := ctx.Get("Authorization")
	config, _ := config.LoadConfig(".")
	secret := config.TokenSecret

	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	tokenString := parts[1]

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	return claims.CustID, nil
}
