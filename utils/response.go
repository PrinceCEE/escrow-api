package utils

type ApiResponse struct {
	Success    *bool  `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Meta       any    `json:"meta,omitempty"`
	StatusCode *int   `json:"-"`
}
