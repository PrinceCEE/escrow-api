package models

type Business struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	ModelMixin
}
