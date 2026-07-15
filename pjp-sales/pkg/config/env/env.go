package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigEnv interface {
	Get(key string) string
}

type configImpl struct {
}

func (config *configImpl) Get(key string) string {
	return os.Getenv(key)
}

func NewCfgEnv(filenames ...string) ConfigEnv {
	err := godotenv.Load(filenames...)
	if err != nil {
		panic(err)
	}
	return &configImpl{}
}

func ValidateRequired(cfg ConfigEnv, keys ...string) error {
	missing := make([]string, 0)
	for _, key := range keys {
		if strings.TrimSpace(cfg.Get(key)) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}
