package config

import (
	"finance/pkg/config/env"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host        string
	User        string
	Password    string
	DBName      string
	DBNumber    int
	Port        int
	DebugMode   bool
	MaxConn     int
	MaxIdle     int
	MaxLifetime int
}

func PostgreSQLConnection(envFile env.ConfigEnv) *gorm.DB {
	var connection *gorm.DB
	var err error
	port, _ := strconv.Atoi(envFile.Get("DB_PORT"))
	debug, err := strconv.ParseBool(envFile.Get("DB_DEBUG"))
	if err != nil {
		panic(err.Error())
	}
	maxConn, _ := strconv.Atoi(envFile.Get("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(envFile.Get("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(envFile.Get("DB_MAX_LIFETIME_CONNECTIONS"))
	config := DBConfig{
		Host:        envFile.Get("DB_HOST"),
		User:        envFile.Get("DB_USER"),
		Password:    envFile.Get("DB_PASS"),
		DBName:      envFile.Get("DB_NAME"),
		Port:        port,
		DebugMode:   debug,
		MaxConn:     maxConn,
		MaxIdle:     maxIdleConn,
		MaxLifetime: maxLifetimeConn,
	}

	var isDebug = logger.Error

	if config.DebugMode {
		isDebug = logger.Info
	}

	dsn := "host=" + config.Host + " port=" + strconv.Itoa(config.Port) + " user=" + config.User + " dbname=" + config.DBName + " password=" + config.Password
	connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(isDebug),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := connection.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(config.MaxConn)
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)

	return connection
}
