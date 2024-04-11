package users

type changePasswordDto struct {
	Password string `json:"password" validate:"required,min=8"`
}

type updateAccountDto struct {
	FirstName    *string `json:"first_name" validate:"omitempty,alpha"`
	LastName     *string `json:"last_name" validate:"omitempty,alpha"`
	BusinessName *string `json:"business_name" validate:"omitempty,alpha"`
	Email        *string `json:"email" validate:"omitempty,email"`
	PhoneNumber  *string `json:"phone_number" validate:"omitempty,min=8"`
	ImageUrl     *string `json:"image_url" validate:"omitempty,url"`
}
