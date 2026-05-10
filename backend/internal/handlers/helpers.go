package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes v as JSON to w with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
