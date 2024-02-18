package auth

import "github.com/Bupher-Co/bupher-api/pkg"

type signUpDto struct {
	AccountType  *string       `json:"account_type" validate:"oneof=personal business"`
	Email        *string       `json:"email" validate:"required,email"`
	PhoneNumber  *string       `json:"phone_number" validate:"min=8"`
	FirstName    *string       `json:"first_name" validate:"alpha"`
	LastName     *string       `json:"last_name" validate:"alpha"`
	Password     *string       `json:"password" validate:"min=8"`
	BusinessName *string       `json:"business_name" validate:"alphanum"`
	RegStage     *pkg.RegStage `json:"reg_stage" validate:"required,numeric,oneof=1 2 3"`
}
