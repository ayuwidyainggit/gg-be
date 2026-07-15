package config

import (
	"strconv"
	"strings"

	"inventory/pkg/config/env"
)

type SAPReplenishmentStatusConfig struct {
	Enabled           bool
	ClientID          string
	SecretKey         string
	TimestampSkewSecs int64 
}

const (
	defaultSAPTimestampSkewSecs = 300
)

func LoadSAPReplenishmentStatusConfig(envCfg env.ConfigEnv) SAPReplenishmentStatusConfig {
	enabled := strings.EqualFold(strings.TrimSpace(envCfg.Get("SAP_REPL_STATUS_ENABLED")), "true") ||
		envCfg.Get("SAP_REPL_STATUS_ENABLED") == "1"

	skew := parseInt64Default(envCfg.Get("SAP_REPL_STATUS_TIMESTAMP_SKEW_SEC"), defaultSAPTimestampSkewSecs)
	if skew <= 0 {
		skew = defaultSAPTimestampSkewSecs
	}

	return SAPReplenishmentStatusConfig{
		Enabled:           enabled,
		ClientID:          strings.TrimSpace(envCfg.Get("SAP_REPL_STATUS_CLIENT_ID")),
		SecretKey:         envCfg.Get("SAP_REPL_STATUS_SECRET_KEY"),
		TimestampSkewSecs: skew,
	}
}

func parseInt64Default(s string, def int64) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return n
}
