package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"mobile/pkg/str"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisConfig represents config for connecting to Redis
type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
	Debug    bool
}

func (r *RedisConfig) SetConfig() *RedisConfig {
	db := str.StringToInt(os.Getenv("REDIS_DB_NUMBER"), 0)
	fmt.Println(os.Getenv("REDIS_DB"), os.Getenv("REDIS_HOST"), "REDISDB")
	r.Addr = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	r.Username = os.Getenv("REDIS_USERNAME")
	r.Password = os.Getenv("REDIS_PASSWORD")
	r.DB = db
	r.Debug, _ = strconv.ParseBool(os.Getenv("REDIS_DEBUG"))
	return r
}

// RedisInstance to instantiation a redis client
func (r *RedisConfig) RedisInstance() (*redis.Client, error) {
	opt := &redis.Options{
		Addr:         r.Addr,
		Password:     r.Password,
		DB:           r.DB,
		Username:     r.Username,
		ReadTimeout:  -1,
		WriteTimeout: -1,
	}
	// logger.Info().Interface("config", r).Msg("redisconfig")
	if !r.Debug {
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	client := redis.NewClient(opt)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

// RedisPing to ping redis server
func (r *RedisConfig) RedisPing() error {
	client, _ := r.RedisInstance()
	if _, err := client.Ping(ctx).Result(); err != nil {
		client.Close()
		r.RedisInstance()
	}

	return nil
}
