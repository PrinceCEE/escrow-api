package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiResponse struct {
	Success    *bool  `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Meta       any    `json:"meta,omitempty"`
	StatusCode *int   `json:"-"`
}

func SendResponse(w http.ResponseWriter, b ApiResponse, headers ...map[string]string) {
	if b.Success == nil {
		*b.Success = true
	}

	if b.StatusCode == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(*b.StatusCode)
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header().Set(k, v)
		}
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	w.Write(jsonData)
}

func SendErrorResponse(w http.ResponseWriter, b ApiResponse, statusCode int) {
	s := false
	b.Success = &s
	b.StatusCode = &statusCode
	SendResponse(w, b)
}
