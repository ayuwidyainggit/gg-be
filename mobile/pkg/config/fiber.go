package config

import (
	"mobile/pkg/config/env"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig(envFile env.ConfigEnv) fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(envFile.Get("SERVER_READ_TIMEOUT"))

	// Parse body limit from env (default: 200MB for multipart form with multiple files)
	bodyLimitMB := 200
	if bodyLimitStr := envFile.Get("BODY_LIMIT_MB"); bodyLimitStr != "" {
		if parsedLimit, err := strconv.Atoi(bodyLimitStr); err == nil && parsedLimit > 0 {
			bodyLimitMB = parsedLimit
		}
	}

	// // Return Fiber configuration.
	return fiber.Config{
		ServerHeader:      envFile.Get("APP_NAME") + " " + envFile.Get("APP_VERSION") + " - " + envFile.Get("APP_STATUS"),
		AppName:           envFile.Get("APP_NAME") + " " + envFile.Get("APP_VERSION") + " - " + envFile.Get("APP_STATUS"),
		StreamRequestBody: false,
		ReadTimeout:       time.Second * time.Duration(readTimeoutSecondsCount),
		EnablePrintRoutes: false,
		BodyLimit:         bodyLimitMB * 1024 * 1024,
	}
}
