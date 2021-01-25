package models

const (
	ErrorTypeBadRequest = "ERROR_BAD_REQUEST"
	ErrorTypeInternalError = "ERROR_INTERNAL"
	ErrorTypeNotFound = "ERROR_NOT_FOUND"
	ErrorTypeValidationError = "ERROR_VALIDATION"
)

const (
	ValidationErrorsRequired = "property is required"
)

// Error when calling the api.
type ApiError struct {
	// The type of error that happened.
	Type string `json:"type,omitempty"`
	// The description of the error.
	Error string `json:"error,omitempty"`
	// A map of data attached to the error. In case of a validation errors this contains the errors per validated property.
	Data map[string]string `json:"data,omitempty"`
}