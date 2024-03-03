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
	BusinessName *string         `json:"business_name" validate:"omitempty,alphanum"`
	RegStage     *utils.RegStage `json:"reg_stage" validate:"required,numeric,oneof=1 2 3"`
}
