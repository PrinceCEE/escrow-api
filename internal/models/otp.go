package models

import (
	"github.com/gofrs/uuid"
)

const (
	SmsOtpType   = "SMS"
	EmailOtpType = "EMAIL"
)

type Otp struct {
	UserID  uuid.UUID `json:"user_id"`
	Code    string    `json:"code"`
	IsUsed  bool      `json:"is_used"`
	OtpType string    `json:"otp_type"`
	ModelMixin
}
