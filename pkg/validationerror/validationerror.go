package validationerror

import "fmt"

type ValidationError struct {
	Message string
	Code    Code
}

// New returns new validation error
func New(message string, code Code) *ValidationError {
	return &ValidationError{
		Message: message,
		Code:    code,
	}
}

// Error implements error interface
func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s (code: %d)", v.Message, v.Code)
}
