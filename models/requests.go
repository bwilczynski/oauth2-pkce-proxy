package models

import (
	"net/http"
	"net/url"
)

const (
	CodeChallengeMethodS256 = "S256"
)

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
	return Validate(
		RequiredField{Name: "client_id", Value: &ar.ClientId},
		RequiredField{Name: "code_challenge", Value: &ar.CodeChallenge},
		RequiredField{Name: "redirect_uri", Value: &ar.RedirectUri},
	)
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
	return Validate(
		RequiredField{Name: "code", Value: &cr.Code},
		RequiredField{Name: "redirect_uri", Value: &cr.RedirectUri},
	)
}

func (cr *AuthorizeCodeRequest) FromQueryParams(r *http.Request) {
	query := r.URL.Query()

	cr.Code = query.Get("code")
	cr.RedirectUri = query.Get("redirect_uri")
}
