package middleware

import (
	"fmt"
	"inventory/pkg/config"
	"inventory/pkg/sap"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	hdrSAPClientID   = "X-Client-Id"
	hdrSAPTimestamp  = "X-Timestamp"
	hdrSAPSignature  = "X-Signature"
	localhostSAPUser = int64(-1)
)

const (
	ErrSAPReplDisabled   = "sap_replenishment_status_disabled"
	localKeySAPUpdatedBy = "sap_integration_user_id"
)

func SAPReplenishmentStatusLocals(c *fiber.Ctx) int64 {
	if v := c.Locals(localKeySAPUpdatedBy); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return localhostSAPUser
}

func SAPReplenishmentStatusProtected(cfg config.SAPReplenishmentStatusConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !cfg.Enabled {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  ErrSAPReplDisabled,
				"message": "SAP replenishment status callback is disabled",
			})
		}

		if cfg.ClientID == "" || cfg.SecretKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "SAP callback is not configured (client id / secret)",
			})
		}

		clientIDHdr := strings.TrimSpace(c.Get(hdrSAPClientID))
		tsStr := strings.TrimSpace(c.Get(hdrSAPTimestamp))
		sigHdr := strings.TrimSpace(c.Get(hdrSAPSignature))

		if clientIDHdr != cfg.ClientID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid X-Client-Id",
			})
		}

		tsInt, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil || tsStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or missing X-Timestamp",
			})
		}

		now := time.Now().Unix()
		skew := cfg.TimestampSkewSecs
		if skew <= 0 {
			skew = 300
		}
		if tsInt > now+skew || tsInt < now-skew {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("X-Timestamp outside allowed skew (±%d seconds)", skew),
			})
		}

		if sigHdr == "" || !sap.ValidMACConstantTime(cfg.SecretKey, clientIDHdr, tsStr, sigHdr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid X-Signature",
			})
		}

		c.Locals(localKeySAPUpdatedBy, localhostSAPUser)
		return c.Next()
	}
}
