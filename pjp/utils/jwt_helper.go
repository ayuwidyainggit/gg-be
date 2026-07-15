package utils

import (
	"strings"

	"github.com/golang-jwt/jwt"
)

// GetUserIdFromToken extracts user_id from JWT token string
func GetUserIdFromToken(token, tokenSecret string) (int64, bool) {
	if token == "" {
		return 0, false
	}

	payload, err := DecodeToken(token, tokenSecret)
	if err != nil {
		return 0, false
	}

	claims, ok := payload.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	userIdRaw, exists := claims["user_id"]
	if !exists {
		return 0, false
	}

	// Handle both float64 (from JSON) and int64
	var userId int64
	switch v := userIdRaw.(type) {
	case float64:
		userId = int64(v)
	case int64:
		userId = v
	case int:
		userId = int64(v)
	default:
		return 0, false
	}

	return userId, true
}

// ExtractTokenFromHeader extracts Bearer token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	fields := strings.Fields(authHeader)
	if len(fields) != 0 && fields[0] == "Bearer" {
		return fields[1]
	}
	return ""
}
