package models

import "github.com/gofrs/uuid"

const (
	PersonalAccountType = "personal"
	BusinessAccountType = "business"
)

type User struct {
	Email                 string     `json:"email,omitempty" db:"email"`
	PhoneNumber           NullString `json:"phone_number,omitempty" db:"phone_number"`
	FirstName             NullString `json:"first_name" db:"first_name"`
	LastName              NullString `json:"last_name" db:"last_name"`
	IsPhoneNumberVerified bool       `json:"is_phone_number_verified" db:"is_phone_number_verified"`
	IsEmailVerified       bool       `json:"is_email_verified" db:"is_email_verified"`
	RegStage              int        `json:"reg_stage" db:"reg_stage"`
	AccountType           string     `json:"account_type" db:"account_type"`
	BusinessID            uuid.UUID  `json:"business_id" db:"business_id"`
	Business              Business   `json:"business,omitempty" db:"-"`
	ImageUrl              string     `json:"image_url,omitempty" db:"image_url"`
	ModelMixin
}
