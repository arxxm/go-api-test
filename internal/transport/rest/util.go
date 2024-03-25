package rest

import (
	"encoding/json"
	"net/http"
)

func response200(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(body)
}

func responseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	body := map[string]string{
		"error": message,
	}
	_ = json.NewEncoder(w).Encode(body)
}

func response201(w http.ResponseWriter, msg interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(msg)
}
