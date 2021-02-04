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

func (h *accessTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		clientID     = "client_id"
		code         = "code"
		codeVerifier = "code_verifier"
	)

	r.ParseForm()
	if err := m.ValidateRequired(r.Form, clientID, code, codeVerifier); err != nil {
		h.logger.Error().Err(err).Msg("")
		writeError(w, err)
		return
	}

	challenge := h.store.Get(r.Form.Get(code))

	cv := m.CodeVerifier(r.Form.Get(codeVerifier))
	if ok := cv.Verify(challenge); !ok {
		h.logger.Error().
			Str("code_verifier", string(cv)).
			Str("challenge", challenge).
			Msg("code verifier not valid")

		http.Error(w, "Code verifier not valid", http.StatusForbidden)
		return
	}

	fd := formData(r.Form)
	fd.Set("client_secret", h.provider.ClientSecret)
	resp, err := http.PostForm(h.provider.TokenEndpoint, fd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func formData(v url.Values) url.Values {
	keys := []string{"client_id", "code", "grant_type"}
	r := make(url.Values)
	for _, k := range keys {
		r.Set(k, v.Get(k))
	}
	return r
}
