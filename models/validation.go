package models

type (
	ValidationError struct {
		FieldName string `json:"fieldName"`
		Message   string `json:"message"`
	}

	ValidationResult struct {
		Message string            `json:"message"`
		Errors  []ValidationError `json:"errors"`
	}

	RequiredField struct {
		Name  string
		Value *string
	}
)

func Validate(fields ...interface{}) (res *ValidationResult, ok bool) {
	errors := []ValidationError{}

	for _, f := range fields {
		if rf, ok := f.(RequiredField); ok && *rf.Value == "" {
			errors = append(errors, ValidationError{FieldName: rf.Name, Message: "Required field"})
		}
	}

	if len(errors) > 0 {
		ok = false
		res = &ValidationResult{Message: "Bad request", Errors: errors}
		return
	}

	ok = true
	return
}
