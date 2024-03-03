package models

import (
	"github.com/gofrs/uuid"
)

const (
	SmsOtpType   = "SMS"
	EmailOtpType = "EMAIL"
)

type Otp struct {
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	Code    string    `json:"code" db:"code"`
	IsUsed  bool      `json:"is_used" db:"is_used"`
	OtpType string    `json:"otp_type" db:"otp_type"`
	ModelMixin
}
