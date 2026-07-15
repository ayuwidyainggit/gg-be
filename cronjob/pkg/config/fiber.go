package config

import (
	"cronjob/pkg/config/env"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewFiberConfig(envFile env.ConfigEnv) fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(envFile.Get("SERVER_READ_TIMEOUT"))

	// Return Fiber configuration.
	return fiber.Config{
		ServerHeader:      envFile.Get("APP_NAME") + " " + envFile.Get("APP_VERSION") + " - " + envFile.Get("APP_STATUS"),
		AppName:           envFile.Get("APP_NAME") + " " + envFile.Get("APP_VERSION") + " - " + envFile.Get("APP_STATUS"),
		StreamRequestBody: false,
		ReadTimeout:       time.Second * time.Duration(readTimeoutSecondsCount),
		EnablePrintRoutes: false,
	}
}
