package env

import (
	"os"

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
