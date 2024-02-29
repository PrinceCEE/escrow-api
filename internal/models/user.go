package models

type User struct {
	Email                 *string `json:"email,omitempty"`
	PhoneNumber           *string `json:"phone_number,omitempty"`
	FirstName             string  `json:"first_name"`
	LastName              string  `json:"last_name"`
	IsPhoneNumberVerified bool    `json:"is_phone_number_verified"`
	IsEmailVerified       bool    `json:"is_email_verified"`
	RegStage              int     `json:"reg_stage"`
	AccountType           string  `json:"account_type"`
	ModelMixin
}
