package errorResponse

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse defines the structure for error responses.
type ErrorResponse struct {
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}

// SendError sends an error response in JSON format.
func SendError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Convert the error message to string (in case it's wrapped or has additional context)
	errorMessage := err.Error()

	// Create and send the JSON response
	response := ErrorResponse{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}

	_ = json.NewEncoder(w).Encode(response) // Ignore encoding errors as we're in an error flow
}
