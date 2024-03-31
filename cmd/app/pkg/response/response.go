package response

import (
	"errors"
	"net/http"

	"github.com/Bupher-Co/bupher-api/pkg/json"
)

var (
	ErrBadRequest      = errors.New("bad request")
	ErrInternalServer  = errors.New("internal server")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrNotFound        = errors.New("not found")
	ErrNotImplemented  = errors.New("not implemented")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrPayloadTooLarge = errors.New("payload too large")
)

type ApiResponseMeta struct {
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	Total        int    `json:"total,omitempty"`
	TotalPages   int    `json:"total_pages,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ApiResponse struct {
	Success    *bool           `json:"success"`
	Message    string          `json:"message"`
	Data       any             `json:"data,omitempty"`
	Meta       ApiResponseMeta `json:"meta,omitempty"`
	StatusCode *int            `json:"-"`
}

func SendResponse(w http.ResponseWriter, b ApiResponse, headers ...map[string]string) {
	if b.Success == nil {
		s := true
		b.Success = &s
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header().Set(k, v)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if b.StatusCode == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(*b.StatusCode)
	}

	jsonData, err := json.WriteJSON(b)
	if err != nil {
		panic(err)
	}

	w.Write(jsonData)
}

func SendErrorResponse(w http.ResponseWriter, a ApiResponse, statusCode int) {
	s := false
	a.Success = &s
	a.StatusCode = &statusCode
	SendResponse(w, a)
}
