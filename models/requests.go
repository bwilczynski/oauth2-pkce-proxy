package models

import (
	"fmt"
	"net/http"
	"strings"
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
	qp := []string{}
	if ar.ClientId != "" {
		qp = append(qp, fmt.Sprintf("client_id=%s", ar.ClientId))
	}
	if ar.CodeChallenge != "" {
		qp = append(qp, fmt.Sprintf("code_challenge=%s", ar.CodeChallenge))
	}
	if ar.CodeChallengeMethod != "" {
		qp = append(qp, fmt.Sprintf("code_challenge_method=%s", ar.CodeChallengeMethod))
	}
	if ar.RedirectUri != "" {
		qp = append(qp, fmt.Sprintf("redirect_uri=%s", ar.RedirectUri))
	}

	return strings.Join(qp, "&")
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

func (cr *AuthorizeCodeRequest) FromQueryParams(r *http.Request) {
	query := r.URL.Query()

	cr.Code = query.Get("code")
	cr.RedirectUri = query.Get("redirect_uri")
}
