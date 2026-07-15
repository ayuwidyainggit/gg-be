package passwordreset

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ResetTokenHMACHex returns a deterministic hex digest for storing and looking up reset tokens.
func ResetTokenHMACHex(secret, plainToken string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(plainToken))
	return hex.EncodeToString(mac.Sum(nil))
}
