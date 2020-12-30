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
	ar := &m.AuthorizeRequest{}
	ar.FromQueryParams(r)

	if err := ar.Validate(); err != nil {
		h.logger.Error().Err(err).Msg("")
		writeError(w, err)
		return
	}

	q := r.URL.Query()
	q.Set("redirect_uri", fmt.Sprintf("http://%s%s?redirect_uri=%s", r.Host, h.callbackPath, ar.RedirectUri))

	redirectURI := fmt.Sprintf("%s?%s", h.provider.AuthorizeURL, q.Encode())
	w.Header().Add("Location", redirectURI)
	setChallengeCookie(w, ar.CodeChallenge)
	w.WriteHeader(http.StatusFound)
}

func setChallengeCookie(w http.ResponseWriter, challenge string) {
	cookie := http.Cookie{Name: challengeCookieName, Value: url.QueryEscape(challenge), Path: "/", HttpOnly: true, MaxAge: challengeCookieMaxAge}
	http.SetCookie(w, &cookie)
}
