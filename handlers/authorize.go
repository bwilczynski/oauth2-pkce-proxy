package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/bwilczynski/oauth2-pkce-proxy/models"
)

type authorizeHandler struct {
	log          *log.Logger
	provider     *models.OAuth2Provider
	callbackPath string
}

func NewAuthorizeHandler(log *log.Logger, provider *models.OAuth2Provider, callbackPath string) *authorizeHandler {
	return &authorizeHandler{log, provider, callbackPath}
}

func (h *authorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar := &models.AuthorizeRequest{}
	ar.FromQueryParams(r)

	h.log.Printf("Authorize handler called: %#v", ar)

	if res, ok := ar.Validate(); !ok {
		h.log.Printf("Bad request: %#v", res)
		writeError(w, res)
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
