package config

import (
	"context"
	"database/sql"
	"fmt"
	"master/pkg/config/env"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/qustavo/sqlhooks/v2"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

// Hooks satisfies the sqlhook.Hooks interface
type Hooks struct{}

var (
	_, currentFile, _, _ = runtime.Caller(0)
	projectRoot          = filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../.."))
)

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	// Replace $n placeholders directly
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		query = strings.Replace(query, placeholder, formatArg(arg), 1)
	}

	// Flatten whitespace
	query = strings.Join(strings.Fields(query), " ")

	if caller := queryCaller(); caller != "" {
		log.Debugf("SQL [%s]: %s", caller, query)
	} else {
		log.Debugf("SQL: %s", query)
	}

	return context.WithValue(ctx, "begin", time.Now()), nil
}

// // After hook will get the timestamp registered on the Before hook and print the elapsed time
// func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
// 	begin := ctx.Value("begin").(time.Time)
// 	// Replace $ placeholders with actual values
// 	for i, arg := range args {
// 		placeholder := fmt.Sprintf("$%d", i+1)
// 		value := fmt.Sprintf("'%v'", arg)
// 		query = strings.Replace(query, placeholder, value, 1)
// 	}

// 	query = str.Replacer(query, strings.NewReplacer("'<nil>'", "null"))
// 	query = str.Replacer(query, strings.NewReplacer(" +0000 UTC", ""))
// 	query = str.Replacer(query, strings.NewReplacer(" +0000", ""))

// 	log.Debugf("> %s < took: %s\n", query, time.Since(begin))
// 	return ctx, nil
// }

func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin, _ := ctx.Value("begin").(time.Time)

	log.Debugf("took: %s", time.Since(begin))

	return ctx, nil
}

func formatArg(arg interface{}) string {
	switch v := arg.(type) {
	case nil:
		return "NULL"
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case []string:
		// for IN (...)
		vals := make([]string, len(v))
		for i, s := range v {
			vals[i] = "'" + strings.ReplaceAll(s, "'", "''") + "'"
		}
		return strings.Join(vals, ",")
	case []int, []int64:
		return strings.Trim(strings.Replace(fmt.Sprint(v), " ", ",", -1), "[]")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func queryCaller() string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	if n == 0 {
		return ""
	}

	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()

		if strings.HasPrefix(frame.File, projectRoot) && frame.File != currentFile {
			if rel, err := filepath.Rel(projectRoot, frame.File); err == nil {
				return fmt.Sprintf("%s:%d", filepath.ToSlash(rel), frame.Line)
			}

			return fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		}

		if !more {
			break
		}
	}

	return ""
}

func ConnToDb(envFile env.ConfigEnv) (sqlx *sqlx.DB, err error) {
	isDebug, err := strconv.ParseBool(envFile.Get("DB_DEBUG"))
	if err != nil {
		panic(err)
	}
	if isDebug {
		sqlx, err = PostgreSQLConnectionHook(envFile)
	} else {
		sqlx, err = PostgreSQLConnection(envFile)
	}
	return
}

func PostgreSQLConnection(envFile env.ConfigEnv) (*sqlx.DB, error) {
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
