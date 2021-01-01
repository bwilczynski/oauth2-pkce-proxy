package handlers

import (
	"io"
	"net/http"
	"net/url"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/rs/zerolog"
)

type accessTokenHandler struct {
	logger   *zerolog.Logger
	store    ChallengeStore
	provider *m.OAuth2Provider
}

func NewAccessTokenHandler(logger *zerolog.Logger, store ChallengeStore, provider *m.OAuth2Provider) *accessTokenHandler {
	return &accessTokenHandler{logger, store, provider}
}

func (h *accessTokenHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tr := &m.AccessTokenRequest{}
	tr.FromQueryParams(r)

	challenge := h.store.Get(tr.Code)
	if ok := tr.CodeVerifier.Verify(challenge); !ok {
		h.logger.Error().
			Str("code_verifier", string(tr.CodeVerifier)).
			Str("challenge", challenge).
			Msg("code verifier not valid")

		http.Error(rw, "Code verifier not valid", http.StatusForbidden)
		return
	}

	fd := formData(r.Form)
	fd.Set("client_secret", h.provider.ClientSecret)
	resp, err := http.PostForm(h.provider.TokenEndpoint, fd)
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
