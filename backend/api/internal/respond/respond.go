package respond

import (
	"encoding/json"
	"net/http"
)

func WithJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		WithError(w, "failed to encode", http.StatusInternalServerError)
	}
}

// Error is a standard error response.
type Error struct {
	// Message is the error message.
	Message string `json:"message" example:"The request was invalid."`
	// StatusCode is the HTTP status code.
	StatusCode int `json:"statusCode" example:"400"`
	// Issues is a list of issues with the request. This is optional.
	Issues []string `json:"issues,omitempty" example:"The request was invalid."`
}

func WithError(w http.ResponseWriter, msg string, status int, issues ...string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	v := Error{
		Message:    msg,
		StatusCode: status,
		Issues:     issues,
	}
	json.NewEncoder(w).Encode(v)
}
