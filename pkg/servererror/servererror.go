package servererror

import (
	"fmt"
	"strings"

	"github.com/corbado/corbado-go/pkg/generated/common"
)

type ServerError struct {
	Details    *string  `json:"details,omitempty"`
	Links      []string `json:"links"`
	Type       string   `json:"type"`
	Validation *[]struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
	Data *map[string]any

	HTTPStatusCode int32              `json:"httpStatusCode"`
	Message        string             `json:"message"`
	RequestData    common.RequestData `json:"requestData"`
	Runtime        float32            `json:"runtime"`
}

// New wraps an error response into a ServerError
func New(cause *common.ErrorRsp) *ServerError {
	if cause == nil {
		return nil
	}

	return &ServerError{
		Details:    cause.Error.Details,
		Links:      cause.Error.Links,
		Type:       cause.Error.Type,
		Validation: cause.Error.Validation,
		Data:       cause.Data,

		HTTPStatusCode: cause.HttpStatusCode,
		Message:        cause.Message,
		RequestData:    cause.RequestData,
		Runtime:        cause.Runtime,
	}
}

// Error implements error interface
func (s *ServerError) Error() string {
	msg := fmt.Sprintf("[%d %s]", s.HTTPStatusCode, s.Message)

	validation := s.GetValidationMessage()
	if validation != "" {
		msg = fmt.Sprintf("%s %s", msg, validation)
	}

	if s.Details != nil && *s.Details != "" {
		msg = fmt.Sprintf("%s (details=%s)", msg, *s.Details)
	}

	msg = fmt.Sprintf("%s (requestID=%s, type=%s)", msg, s.RequestData.RequestID, s.Type)

	return msg
}

// GetValidationMessage returns all validation messages as one string
func (s *ServerError) GetValidationMessage() string {
	if s.Validation == nil || len(*s.Validation) == 0 {
		return ""
	}

	fieldMessages := make([]string, len(*s.Validation))
	for i, validation := range *s.Validation {
		fieldMessages[i] = fmt.Sprintf("%s: %s", validation.Field, validation.Message)
	}

	return strings.Join(fieldMessages, "; ")
}
