package handlers

import (
	"net/http"
	"net/url"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/rs/zerolog"
)

type authorizeCodeHandler struct {
	logger *zerolog.Logger
	store  ChallengeStore
}

func NewAuthorizeCodeHandler(logger *zerolog.Logger, store ChallengeStore) *authorizeCodeHandler {
	return &authorizeCodeHandler{logger, store}
}

func (h *authorizeCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		code        = "code"
		redirectURI = "redirect_uri"
	)

	q := r.URL.Query()
	if err := m.ValidateRequired(q, code, redirectURI); err != nil {
		h.logger.Error().Err(err).Msg("")
		writeError(w, err)
		return
	}

	challenge := readChallengeCookie(r)
	h.store.Write(q.Get(code), challenge)

	rURI, _ := url.Parse(q.Get(redirectURI))
	q.Del(redirectURI)
	rURI.RawQuery = q.Encode()

	w.Header().Add("Location", rURI.String())
	w.WriteHeader(http.StatusFound)
}

func readChallengeCookie(r *http.Request) string {
	c, err := r.Cookie(challengeCookieName)
	if err != nil {
		return ""
	}
	challenge, _ := url.QueryUnescape(c.Value)
	return challenge
}
