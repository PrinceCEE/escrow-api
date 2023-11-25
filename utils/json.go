package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ReadJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(dst)

	if err != nil {
		return err
	}

	return nil
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
