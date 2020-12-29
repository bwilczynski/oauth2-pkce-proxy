package handlers

import (
	"log"
	"net/http"
	"net/url"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
)

type authorizeCodeHandler struct {
	log   *log.Logger
	store ChallengeStore
}

func NewAuthorizeCodeHandler(log *log.Logger, store ChallengeStore) *authorizeCodeHandler {
	return &authorizeCodeHandler{log, store}
}

func (h *authorizeCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cr := &m.AuthorizeCodeRequest{}
	cr.FromQueryParams(r)

	if err := cr.Validate(); err != nil {
		h.log.Printf("Bad request: %#v", err)
		writeError(w, err)
		return
	}

	cc := readChallengeCookie(r)

	h.store.Write(cr.Code, m.CodeVerifier(cc))

	h.log.Printf("AuthorizeCode handler called: %#v, cc: %s", cr, cc)

	q := r.URL.Query()
	q.Del("redirect_uri")

	redirectURI, _ := url.Parse(cr.RedirectUri)
	redirectURI.RawQuery = q.Encode()

	w.Header().Add("Location", redirectURI.String())
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
