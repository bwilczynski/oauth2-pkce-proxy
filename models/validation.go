package models

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

func Validate(fields ...interface{}) error {
	errors := []FieldError{}

	for _, f := range fields {
		if rf, ok := f.(RequiredField); ok && *rf.Value == "" {
			errors = append(errors, FieldError{FieldName: rf.Name, Message: "Required field"})
		}
	}

	if len(errors) > 0 {
		return &ValidationError{Message: "Bad request", Errors: errors}
	}

	return nil
}
