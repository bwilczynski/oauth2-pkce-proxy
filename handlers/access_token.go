package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
)

type accessTokenHandler struct {
	log      *log.Logger
	store    ChallengeStore
	provider *m.OAuth2Provider
}

func NewAccessTokenHandler(log *log.Logger, store ChallengeStore, provider *m.OAuth2Provider) *accessTokenHandler {
	return &accessTokenHandler{log, store, provider}
}

func (h *accessTokenHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tr := &m.AccessTokenRequest{}
	tr.FromQueryParams(r)

	h.log.Printf("AccessToken handler called: %#v", tr)

	challenge := h.store.Get(tr.Code)
	if ok := tr.CodeVerifier.Verify(challenge); !ok {
		h.log.Printf("Code verifier not valid: %v for challenge: %v", tr.CodeVerifier, challenge)
		http.Error(rw, "Code verifier not valid", http.StatusForbidden)
		return
	}

	fd := formData(r.Form)
	fd.Set("client_secret", h.provider.ClientSecret)
	resp, err := http.PostForm(h.provider.TokenURL.String(), fd)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(rw, resp.Body)
}

func formData(v url.Values) url.Values {
	keys := []string{"client_id", "code", "grant_type"}
	r := make(url.Values)
	for _, k := range keys {
		r.Set(k, v.Get(k))
	}
	return r
}
