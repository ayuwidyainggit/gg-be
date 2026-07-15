package sap

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func SignatureStringToSign(clientID, timestamp string) string {
	return strings.TrimSpace(clientID) + ":" + strings.TrimSpace(timestamp)
}

func HMACSHA256Hex(secret, stringToSign string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(stringToSign))
	return hex.EncodeToString(mac.Sum(nil))
}

func ValidMACConstantTime(secret, clientID, timestamp, suppliedMACHex string) bool {
	expected := HMACSHA256Hex(secret, SignatureStringToSign(clientID, timestamp))
	if len(suppliedMACHex) != len(expected) {
		return false
	}
	in := strings.ToLower(strings.TrimSpace(suppliedMACHex))
	return hmac.Equal([]byte(expected), []byte(in))
}
