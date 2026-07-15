package model

import (
	"time"
)

type PasswordResetRequest struct {
	ID                  int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID              int64      `gorm:"column:user_id" json:"user_id"`
	Email               string     `gorm:"column:email" json:"email"`
	OtpCode             string     `gorm:"column:otp_code" json:"-"`
	OtpExpiredAt        time.Time  `gorm:"column:otp_expired_at" json:"otp_expired_at"`
	OtpAttemptCount     int        `gorm:"column:otp_attempt_count" json:"otp_attempt_count"`
	OtpMaxAttempt       int        `gorm:"column:otp_max_attempt" json:"otp_max_attempt"`
	ResendCount         int        `gorm:"column:resend_count" json:"resend_count"`
	ResendMax           int        `gorm:"column:resend_max" json:"resend_max"`
	ResendCooldownUntil *time.Time `gorm:"column:resend_cooldown_until" json:"resend_cooldown_until"`
	RequestID           string     `gorm:"column:request_id" json:"request_id"`
	ResetToken          *string    `gorm:"column:reset_token" json:"-"`
	ResetTokenExpiredAt *time.Time `gorm:"column:reset_token_expired_at" json:"reset_token_expired_at"`
	Status              string     `gorm:"column:status" json:"status"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (PasswordResetRequest) TableName() string {
	return "sys.password_reset_requests"
}
