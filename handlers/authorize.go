package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/rs/zerolog"
)

type authorizeHandler struct {
	logger       *zerolog.Logger
	provider     *m.OAuth2Provider
	callbackPath string
}

func NewAuthorizeHandler(logger *zerolog.Logger, provider *m.OAuth2Provider, callbackPath string) *authorizeHandler {
	return &authorizeHandler{logger, provider, callbackPath}
}

func (h *authorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		redirectURI   = "redirect_uri"
		codeChallenge = "code_challenge"
	)

	q := r.URL.Query()
	if err := m.ValidateRequired(q, redirectURI, codeChallenge); err != nil {
		h.logger.Error().Err(err).Msg("")
		writeError(w, err)
		return
	}

	q.Set(redirectURI, fmt.Sprintf("http://%s%s?redirect_uri=%s", r.Host, h.callbackPath, q.Get(redirectURI)))
	loc := fmt.Sprintf("%s?%s", h.provider.AuthorizationEndpoint, q.Encode())
	w.Header().Add("Location", loc)
	setChallengeCookie(w, q.Get(codeChallenge))
	w.WriteHeader(http.StatusFound)
}

func setChallengeCookie(w http.ResponseWriter, challenge string) {
	cookie := http.Cookie{Name: challengeCookieName, Value: url.QueryEscape(challenge), Path: "/", HttpOnly: true, MaxAge: challengeCookieMaxAge}
	http.SetCookie(w, &cookie)
}
