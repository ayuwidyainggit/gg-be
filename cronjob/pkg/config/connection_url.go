package config

import (
	"cronjob/pkg/config/env"
	"fmt"
)

// ConnectionURL func for building URL connection.
func ConnectionURL(n string, envCfg env.ConfigEnv) (string, error) {
	// Define URL to connection.
	var url string

	// Switch given names.
	switch n {
	case "mysql":
		url = fmt.Sprintf(
			"%s:%s@(%s:%s)/%s",
			envCfg.Get("DB_USER"),
			envCfg.Get("DB_PASSWORD"),
			envCfg.Get("DB_HOST"),
			envCfg.Get("DB_PORT"),
			envCfg.Get("DB_NAME"),
		)
	case "postgres":
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			envCfg.Get("DB_HOST"),
			envCfg.Get("DB_PORT"),
			envCfg.Get("DB_USER"),
			envCfg.Get("DB_PASSWORD"),
			envCfg.Get("DB_NAME"),
			envCfg.Get("DB_SSL_MODE"),
		)
	case "redis":
		// URL for Redis connection.
		url = fmt.Sprintf(
			"%s:%s",
			envCfg.Get("REDIS_HOST"),
			envCfg.Get("REDIS_PORT"),
		)
	case "fiber":
		// URL for Fiber connection.
		url = fmt.Sprintf(
			"%s:%s",
			envCfg.Get("SERVER_HOST"),
			envCfg.Get("SERVER_PORT"),
		)
	default:
		// Return error message.
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	// Return connection URL.
	return url, nil
}
