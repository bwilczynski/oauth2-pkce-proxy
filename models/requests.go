package models

import (
	"net/http"
	"net/url"
)

const (
	CodeChallengeMethodS256 = "S256"
)

type ValidationError struct {
	FieldName string `json:"fieldName"`
	Message   string `json:"message"`
}

type ValidationResult struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

type AuthorizeRequest struct {
	ClientId            string
	CodeChallenge       string
	CodeChallengeMethod string
	RedirectUri         string
}

func (ar *AuthorizeRequest) FromQueryParams(r *http.Request) {
	query := r.URL.Query()

	ar.ClientId = query.Get("client_id")
	ar.CodeChallenge = query.Get("code_challenge")
	if codeChallengeMethod := query.Get("code_challenge_method"); codeChallengeMethod != "" {
		ar.CodeChallengeMethod = codeChallengeMethod
	} else {
		ar.CodeChallengeMethod = CodeChallengeMethodS256
	}
	ar.RedirectUri = query.Get("redirect_uri")
}

func (ar *AuthorizeRequest) URLQuery() string {
	q := url.Values{}

	if ar.ClientId != "" {
		q.Add("client_id", ar.ClientId)
	}
	if ar.CodeChallenge != "" {
		q.Add("code_challenge", ar.CodeChallenge)
	}
	if ar.CodeChallengeMethod != "" {
		q.Add("code_challenge_method", ar.CodeChallengeMethod)
	}
	if ar.RedirectUri != "" {
		q.Add("redirect_uri", ar.RedirectUri)
	}

	return q.Encode()
}

func (ar *AuthorizeRequest) Validate() (res *ValidationResult, ok bool) {
	errors := []ValidationError{}

	if ar.ClientId == "" {
		errors = append(errors, ValidationError{FieldName: "client_id", Message: "Required field"})
	}
	if ar.CodeChallenge == "" {
		errors = append(errors, ValidationError{FieldName: "code_challenge", Message: "Required field"})
	}
	if ar.RedirectUri == "" {
		errors = append(errors, ValidationError{FieldName: "redirect_uri", Message: "Required field"})
	}

	if len(errors) > 0 {
		ok = false
		res = &ValidationResult{Message: "Bad request", Errors: errors}
		return
	}

	ok = true
	return
}

type AuthorizeCodeRequest struct {
	Code        string
	RedirectUri string
}

func (ar *AuthorizeCodeRequest) URLQuery() string {
	q := url.Values{}
	if ar.Code != "" {
		q.Add("code", ar.Code)
	}
	if ar.RedirectUri != "" {
		q.Add("redirect_uri", ar.RedirectUri)
	}

	return q.Encode()
}

func (cr *AuthorizeCodeRequest) Validate() (res *ValidationResult, ok bool) {
	errors := []ValidationError{}

	if cr.Code == "" {
		errors = append(errors, ValidationError{FieldName: "code", Message: "Required field"})
	}
	if cr.RedirectUri == "" {
		errors = append(errors, ValidationError{FieldName: "redirect_uri", Message: "Required field"})
	}

	if len(errors) > 0 {
		ok = false
		res = &ValidationResult{Message: "Bad request", Errors: errors}
		return
	}

	ok = true
	return
}

func (cr *AuthorizeCodeRequest) FromQueryParams(r *http.Request) {
	query := r.URL.Query()

	cr.Code = query.Get("code")
	cr.RedirectUri = query.Get("redirect_uri")
}
