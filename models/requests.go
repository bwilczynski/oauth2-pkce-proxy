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
	q := make(url.Values)

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

func (ar *AuthorizeRequest) Validate() error {
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
	q := make(url.Values)
	if ar.Code != "" {
		q.Add("code", ar.Code)
	}
	if ar.RedirectUri != "" {
		q.Add("redirect_uri", ar.RedirectUri)
	}

	return q.Encode()
}

func (cr *AuthorizeCodeRequest) Validate() error {
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

type AccessTokenRequest struct {
	ClientId     string
	Code         string
	GrantType    string
	CodeVerifier string
}

func (tr *AccessTokenRequest) FromQueryParams(r *http.Request) {
	r.ParseForm()

	tr.ClientId = r.Form.Get("client_id")
	tr.Code = r.Form.Get("code")
	tr.GrantType = r.Form.Get("grant_type")
	tr.CodeVerifier = r.Form.Get("code_verifier")
}
