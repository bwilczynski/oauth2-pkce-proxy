package models

import (
	"fmt"
	"net/url"
	"strings"
)

type (
	ValidationError struct {
		Message       string   `json:"message"`
		MissingFields []string `json:"missingFields"`
	}
)

func (vr *ValidationError) Error() string {
	return fmt.Sprintf("%s. Missing fields: %s.", vr.Message, strings.Join(vr.MissingFields, ", "))
}

func ValidateRequired(v url.Values, fields ...string) error {
	mf := []string{}
	for _, f := range fields {
		if v.Get(f) == "" {
			mf = append(mf, f)
		}
	}
	if len(mf) > 0 {
		return &ValidationError{Message: "Bad request", MissingFields: mf}
	}
	return nil
}
