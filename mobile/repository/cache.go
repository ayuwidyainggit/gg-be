package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client
}
type CacheRepository interface {
	SaveOTP(ctx context.Context, userID int64, otp string) error
	GetOTP(ctx context.Context, userID int64) (*string, error)
	DeleteOTP(ctx context.Context, userID int64) error
}

func OtpKeyGenerate(userID int64) string {
	return fmt.Sprintf("otp:userid:%v", userID)
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) SaveOTP(ctx context.Context, userID int64, otp string) error {
	fmt.Println("set cache")

	if _, err := c.client.HSet(ctx, OtpKeyGenerate(userID), "otp", otp).Result(); err != nil {
		return fmt.Errorf("SaveOTP: redis error: %w", err)
	}
	c.client.Expire(ctx, OtpKeyGenerate(userID), 60*time.Minute)
	return nil
}

func (c *Cache) GetOTP(ctx context.Context, userID int64) (*string, error) {
	result, err := c.client.HGet(ctx, OtpKeyGenerate(userID), "otp").Result()
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("find: redis error: %w", err)
	}
	if result == "" {
		return nil, fmt.Errorf("OTP not found or expired")
	}
	return &result, nil
}
func (c *Cache) DeleteOTP(ctx context.Context, userID int64) error {
	_, err := c.GetOTP(ctx, userID)
	if err != nil {
		return err
	}
	c.client.Del(ctx, OtpKeyGenerate(userID))
	return nil
}
