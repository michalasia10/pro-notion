package httpx

import "net/http"

// HTTPError represents an error with associated HTTP status code and optional details.
type HTTPError struct {
	StatusCode int
	Message    string
	Details    any
}

func (e *HTTPError) Error() string { return e.Message }

func BadRequest(msg string, details any) *HTTPError {
	return &HTTPError{StatusCode: http.StatusBadRequest, Message: msg, Details: details}
}

func NotFound(msg string) *HTTPError {
	return &HTTPError{StatusCode: http.StatusNotFound, Message: msg}
}

func Unprocessable(msg string, details any) *HTTPError {
	return &HTTPError{StatusCode: http.StatusUnprocessableEntity, Message: msg, Details: details}
}
