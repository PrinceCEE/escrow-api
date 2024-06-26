package models

import (
	"time"
)

const (
	SmsOtpType        = "sms"
	EmailOtpType      = "email"
	ResetPasswordType = "reset_password"
	OtpExpiresIn      = 10 // 10 minutes
)

type Otp struct {
	UserID    string    `json:"user_id" db:"user_id"`
	Code      string    `json:"code" db:"code"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
	OtpType   string    `json:"otp_type" db:"otp_type"`
	ExpiresIn time.Time `json:"expires_in" db:"expires_in"`
	ModelMixin
}
