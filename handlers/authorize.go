package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
)

type authorizeHandler struct {
	log          *log.Logger
	provider     *m.OAuth2Provider
	callbackPath string
}

func NewAuthorizeHandler(log *log.Logger, provider *m.OAuth2Provider, callbackPath string) *authorizeHandler {
	return &authorizeHandler{log, provider, callbackPath}
}

func (h *authorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar := &m.AuthorizeRequest{}
	ar.FromQueryParams(r)

	h.log.Printf("Authorize handler called: %#v", ar)

	if err := ar.Validate(); err != nil {
		h.log.Printf("Bad request: %#v", err)
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
