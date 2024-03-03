package models

type User struct {
	Email                 *string `json:"email,omitempty" db:"email"`
	PhoneNumber           *string `json:"phone_number,omitempty" db:"phone_number"`
	FirstName             string  `json:"first_name" db:"first_name"`
	LastName              string  `json:"last_name" db:"last_name"`
	IsPhoneNumberVerified bool    `json:"is_phone_number_verified" db:"is_phone_number_verified"`
	IsEmailVerified       bool    `json:"is_email_verified" db:"is_email_verified"`
	RegStage              int     `json:"reg_stage" db:"reg_stage"`
	AccountType           string  `json:"account_type" db:"account_type"`
	ModelMixin
}
