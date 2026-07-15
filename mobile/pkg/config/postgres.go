package config

import (
	"context"

	"fmt"
	"mobile/pkg/config/env"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/qustavo/sqlhooks/v2"
)

type Hooks struct{}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	fmt.Printf("> %s %q", query, args)
	return context.WithValue(ctx, "begin", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	fmt.Printf(". took: %s\n", time.Since(begin))
	return ctx, nil
}

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

func ConnToDb(envFile env.ConfigEnv) (sqlx *sqlx.DB, err error) {
	isDebug, err := strconv.ParseBool(envFile.Get("DB_DEBUG"))
	if err != nil {
		panic(err)
	}
	if isDebug {
		sqlx, err = PostgreSQLConnectionHook(envFile)
	} else {
		sqlx, err = PostgreSQLConnection2(envFile)
	}
	return
}

func PostgreSQLConnection2(envFile env.ConfigEnv) (*sqlx.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(envFile.Get("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(envFile.Get("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(envFile.Get("DB_MAX_LIFETIME_CONNECTIONS"))

	dbServerUrl := "host= " + envFile.Get("DB_HOST") + " port=" + envFile.Get("DB_PORT") + " user=" + envFile.Get("DB_USER") + " password=" + envFile.Get("DB_PASS") + " dbname=" + envFile.Get("DB_NAME") + " sslmode=disable"
	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", dbServerUrl)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}
	// Set database connection settings.
	db.SetMaxOpenConns(maxConn)                           // the default is 0 (unlimited)
	db.SetMaxIdleConns(maxIdleConn)                       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, connections are reused forever

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil

}

func PostgreSQLConnectionHook(envFile env.ConfigEnv) (*sqlx.DB, error) {

	// Define database connection settings.
	maxConn, _ := strconv.Atoi(envFile.Get("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(envFile.Get("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(envFile.Get("DB_MAX_LIFETIME_CONNECTIONS"))

	sql.Register("pgWithHooks", sqlhooks.Wrap(&pq.Driver{}, &Hooks{}))

	dbServerUrl := "host= " + envFile.Get("DB_HOST") + " port=" + envFile.Get("DB_PORT") + " user=" + envFile.Get("DB_USER") + " password=" + envFile.Get("DB_PASS") + " dbname=" + envFile.Get("DB_NAME") + " sslmode=disable"
	db, err := sql.Open("pgWithHooks", dbServerUrl)
	if err != nil {
		return nil, err
	}

	// pass it sqlx
	sqlxDB := sqlx.NewDb(db, "pgx")
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}
	// Set database connection settings.
	sqlxDB.SetMaxOpenConns(maxConn)                           // the default is 0 (unlimited)
	sqlxDB.SetMaxIdleConns(maxIdleConn)                       // defaultMaxIdleConns = 2
	sqlxDB.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, connections are reused forever

	// Try to ping database.
	if err := sqlxDB.Ping(); err != nil {
		defer sqlxDB.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return sqlxDB, nil
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
