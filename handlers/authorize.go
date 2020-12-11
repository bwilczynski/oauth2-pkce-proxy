package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type OAuth2Provider struct {
	AuthorizeURL *url.URL
	TokenURL     *url.URL
}

func NewOAuth2Provider(authorizeURL *url.URL, tokenURL *url.URL) *OAuth2Provider {
	return &OAuth2Provider{authorizeURL, tokenURL}
}

type Authorize struct {
	log      *log.Logger
	provider *OAuth2Provider
}

func NewAuthorize(log *log.Logger, provider *OAuth2Provider) *Authorize {
	return &Authorize{log, provider}
}

type authorizeRequest struct {
	ClientId            string
	CodeChallenge       string
	CodeChallengeMethod string
	RedirectUri         string
}

const (
	CodeChallengeMethodS256 = "S256"
)

func (ar *authorizeRequest) FromQueryParams(r *http.Request) {
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

func (h *Authorize) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar := &authorizeRequest{}
	ar.FromQueryParams(r)

	h.log.Printf("Authorize handler called: %#v", ar)

	if ar.CodeChallengeMethod != CodeChallengeMethodS256 {
		h.log.Printf("Code challenge method %v not supported", ar.CodeChallengeMethod)
		http.Error(w, "Code challenge method not supported", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	q.Set("redirect_uri", fmt.Sprintf("http://%s", r.Host))

	redirectURI := fmt.Sprintf("%s?%s", h.provider.AuthorizeURL, q.Encode())
	w.Header().Add("Location", redirectURI)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
