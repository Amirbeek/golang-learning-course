package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// Limit the request body size to 4 MB to prevent DoS attacks
	maxBytes := 1024 * 1024 * 4
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a JSON decoder for the request body
	decoder := json.NewDecoder(r.Body)

	// Reject any unknown fields that are not in the struct
	decoder.DisallowUnknownFields()

	// Decode the JSON into the provided struct
	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, envelope{message})
}
