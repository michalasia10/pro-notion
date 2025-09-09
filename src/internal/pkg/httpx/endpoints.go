package httpx

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

// EndpointFunc is a function that handles HTTP requests and returns status, response data, and error
type EndpointFunc func(*http.Request) (int, any, error)

// EndpointJSONFunc is a function that handles HTTP requests with JSON body
type EndpointJSONFunc[T any] func(*http.Request, T) (int, any, error)

// Endpoint creates an HTTP handler from an EndpointFunc
func Endpoint(fn EndpointFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, data, err := fn(r)

		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, status, data)
	}
}

// EndpointJSON creates an HTTP handler from an EndpointJSONFunc that expects a JSON body
func EndpointJSON[T any](fn EndpointJSONFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handleError(w, BadRequest("Invalid JSON body", nil))
			return
		}

		status, data, err := fn(r, body)

		if err != nil {
			handleError(w, err)
			return
		}

		WriteJSON(w, status, data)
	}
}

// handleError handles errors returned from endpoint functions
func handleError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(*HTTPError); ok {
		WriteJSON(w, httpErr.StatusCode, map[string]any{
			"error":   httpErr.Message,
			"details": httpErr.Details,
		})
		return
	}

	// Default internal server error
	WriteJSON(w, http.StatusInternalServerError, map[string]any{
		"error": err.Error(),
	})
}

// ValidateTags validates struct fields with basic validation tags
func ValidateTags(v any) error {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}

		tags := strings.Split(tag, ",")
		fieldName := fieldType.Name

		for _, t := range tags {
			t = strings.TrimSpace(t)

			switch {
			case t == "required":
				if isEmptyValue(field) {
					return BadRequest("Validation failed", map[string]string{
						fieldName: "field is required",
					})
				}
			case strings.HasPrefix(t, "min="):
				// Basic min validation for strings
				if field.Kind() == reflect.String {
					minStr := strings.TrimPrefix(t, "min=")
					if len(minStr) > 0 && len(field.String()) == 0 {
						return BadRequest("Validation failed", map[string]string{
							fieldName: "field cannot be empty",
						})
					}
				}
			case t == "email":
				// Basic email validation
				if field.Kind() == reflect.String {
					email := field.String()
					if email != "" && !strings.Contains(email, "@") {
						return BadRequest("Validation failed", map[string]string{
							fieldName: "invalid email format",
						})
					}
				}
			}
		}
	}

	return nil
}

// isEmptyValue checks if a reflect.Value is empty
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}
