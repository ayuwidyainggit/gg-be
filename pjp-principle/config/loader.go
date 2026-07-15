package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost      string `mapstructure:"POSTGRES_HOST"`
	DBUsername  string `mapstructure:"POSTGRES_USER"`
	DBPassword  string `mapstructure:"POSTGRES_PASSWORD"`
	DBName      string `mapstructure:"POSTGRES_DB"`
	DBPort      string `mapstructure:"POSTGRES_PORT"`
	ServerPort  string `mapstructure:"PORT"`
	SwaggerHost string `mapstructure:"SWAGGER_HOST"`
	SwaggerUrl  string `mapstructure:"SWAGGER_URL"`
	TokenSecret string `mapstructure:"JWT_SECRET_KEY"`
	KongUrl     string `mapstructure:"KONG_URL"`
	EmailFrom   string `mapstructure:"EMAIL_FROM"`
	SMTPHost    string `mapstructure:"SMTP_HOST"`
	SMTPPass    string `mapstructure:"SMTP_PASS"`
	SMTPPort    int    `mapstructure:"SMTP_PORT"`
	SMTPUser    string `mapstructure:"SMTP_USER"`

	Environment string `mapstructure:"ENVIRONMENT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// viper.SetConfigType("env")
	// viper.AddConfigPath(".")
	// viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.SetConfigName("development")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
