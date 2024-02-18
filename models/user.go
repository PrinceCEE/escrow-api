package models

type User struct {
	Email                 string `json:"email"`
	PhoneNumber           string `json:"phone_number"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	IsPhoneNumberVerified bool   `json:"is_phone_number_verified"`
	IsEmailVerified       bool   `json:"is_email_verified"`
	RegStage              int    `json:"reg_stage"`
	ModelMixin
}
