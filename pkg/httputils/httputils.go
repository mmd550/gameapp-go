package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(ErrorResponse{Error: message})

	if err != nil {
		fmt.Printf("Failed to json marshal error message: %v", err)
	}

	_, err = w.Write(jsonError)

	if err != nil {
		fmt.Printf("Failed to write error response: %v", err)
	}
}
