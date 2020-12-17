package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/bwilczynski/oauth2-pkce-proxy/models"
)

type authorizeCodeHandler struct {
	log *log.Logger
}

func NewAuthorizeCodeHandler(log *log.Logger) *authorizeCodeHandler {
	return &authorizeCodeHandler{log}
}

func (h *authorizeCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cr := &models.AuthorizeCodeRequest{}
	cr.FromQueryParams(r)

	if res, ok := cr.Validate(); !ok {
		h.log.Printf("Bad request: %#v", res)
		writeError(w, res)
		return
	}

	cc := readChallengeCookie(r)
	h.log.Printf("AuthorizeCode handler called: %#v, cc: %s", cr, cc)

	q := r.URL.Query()
	q.Del("redirect_uri")

	redirectURI := cr.RedirectUri
	if len(q) > 0 {
		redirectURI += "&" + q.Encode()
	}

	w.Header().Add("Location", redirectURI)
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
