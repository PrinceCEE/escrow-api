package models

type Business struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	ModelMixin
}
