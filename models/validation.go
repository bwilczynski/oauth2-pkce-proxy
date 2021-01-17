package models

import (
	"net/url"
)

type (
	FieldError struct {
		FieldName string `json:"fieldName"`
		Message   string `json:"message"`
	}

	ValidationError struct {
		Message string       `json:"message"`
		Errors  []FieldError `json:"errors"`
	}

	RequiredField struct {
		Name  string
		Value *string
	}
)

func (vr *ValidationError) Error() string {
	return vr.Message
}

func ValidateRequired(v url.Values, fields ...string) error {
	errors := []FieldError{}
	for _, f := range fields {
		if v.Get(f) == "" {
			errors = append(errors, FieldError{FieldName: f, Message: "Required field"})
		}
	}
	if len(errors) > 0 {
		return &ValidationError{Message: "Bad request", Errors: errors}
	}
	return nil
}
