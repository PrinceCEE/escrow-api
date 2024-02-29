package models

const (
	SmsOtpType   = "SMS"
	EmailOtpType = "EMAIL"
)

type Otp struct {
	UserID  string `json:"user_id"`
	Code    string `json:"code"`
	IsUsed  bool   `json:"is_used"`
	OtpType string `json:"otp_type"`
	ModelMixin
}
