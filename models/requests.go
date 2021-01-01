package models

import (
	"net/http"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

const (
	CodeChallengeMethodS256 = "S256"
)

type AuthorizeRequest struct {
	ClientId            string `mapstructure:"client_id"`
	CodeChallenge       string `mapstructure:"code_challenge"`
	CodeChallengeMethod string `mapstructure:"code_challenge_method"`
	RedirectUri         string `mapstructure:"redirect_uri"`
}

func (ar *AuthorizeRequest) FromQueryParams(r *http.Request) {
	decode(r.URL.Query(), ar)

	if ar.CodeChallengeMethod == "" {
		ar.CodeChallengeMethod = CodeChallengeMethodS256
	}
}

func (ar *AuthorizeRequest) URLQuery() string {
	q := encode(ar)
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
	Code        string `mapstructure:"code"`
	RedirectUri string `mapstructure:"redirect_uri"`
}

func (cr *AuthorizeCodeRequest) FromQueryParams(r *http.Request) error {
	return decode(r.URL.Query(), cr)
}

func (ar *AuthorizeCodeRequest) URLQuery() string {
	q := encode(ar)
	return q.Encode()
}

func (cr *AuthorizeCodeRequest) Validate() error {
	return Validate(
		RequiredField{Name: "code", Value: &cr.Code},
		RequiredField{Name: "redirect_uri", Value: &cr.RedirectUri},
	)
}

type AccessTokenRequest struct {
	ClientId     string       `mapstructure:"client_id"`
	Code         string       `mapstructure:"code"`
	GrantType    string       `mapstructure:"grant_type"`
	CodeVerifier CodeVerifier `mapstructure:"code_verifier"`
}

func (tr *AccessTokenRequest) FromQueryParams(r *http.Request) error {
	r.ParseForm()
	return decode(r.Form, tr)
}

func decode(vals url.Values, v interface{}) error {
	fvals := make(map[string]string)

	for k, v := range vals {
		fvals[k] = v[0]
	}

	return mapstructure.Decode(fvals, v)
}

func encode(v interface{}) url.Values {
	vals := make(url.Values)
	fvals := make(map[string]string)

	mapstructure.Decode(v, &fvals)
	for k, v := range fvals {
		vals.Set(k, v)
	}

	return vals
}
