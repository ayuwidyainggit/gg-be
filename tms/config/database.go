package config

import (
	"fmt"
	"scyllax-tms/exception"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectionDB(config *Config) *gorm.DB {

	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, config.DBName, config.DBTimeZone)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}
	//helper.ErrorPanic(err)

	connection, err := db.DB()
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(100)
	connection.SetConnMaxLifetime(time.Second * time.Duration(300))

	fmt.Println("🚀 Connected Successfully to the Database")
	return db
}
