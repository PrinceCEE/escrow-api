package auth

import (
	"github.com/Bupher-Co/bupher-api/pkg/utils"
)

type signUpDto struct {
	AccountType  *string         `json:"account_type" validate:"omitempty,oneof=personal business"`
	Email        *string         `json:"email" validate:"required,email"`
	PhoneNumber  *string         `json:"phone_number" validate:"omitempty,min=8"`
	FirstName    *string         `json:"first_name" validate:"omitempty,alpha"`
	LastName     *string         `json:"last_name" validate:"omitempty,alpha"`
	Password     *string         `json:"password" validate:"omitempty,min=8"`
	BusinessName *string         `json:"business_name" validate:"omitempty"`
	RegStage     *utils.RegStage `json:"reg_stage" validate:"required,numeric,oneof=1 2 3"`
}

type verifyCodeDto struct {
	Email   string `json:"email" validate:"required"`
	Code    string `json:"code" validate:"required,len=4"`
	OtpType string `json:"otp_type" validate:"required,oneof=sms email reset_password"`
}

type signInDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type forgotPasswordDto struct {
	Email string `json:"email" validate:"required,email"`
}

type changePasswordDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type resendCodeOTPDto struct {
	Identifier string `json:"identifier" validate:"required,email"`
	OtpType string `json:"otp_type" validate:"required,oneof=sms email reset_password"`
}