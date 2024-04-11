package models

type Business struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	ImageUrl string `json:"image_url,omitempty" db:"image_url"`
	ModelMixin
}
