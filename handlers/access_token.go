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

	// TODO: validate received code challenge with code verifier
	// h.store.Get(code)

	fd := formData(r.Form)
	fd.Set("client_secret", h.provider.ClientSecret)
	resp, err := http.PostForm(h.provider.TokenURL.String(), fd)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
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
