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
	cr := &m.AuthorizeCodeRequest{}
	cr.FromQueryParams(r)

	if err := cr.Validate(); err != nil {
		h.logger.Error().Err(err).Msg("")
		writeError(w, err)
		return
	}

	challenge := readChallengeCookie(r)
	h.store.Write(cr.Code, challenge)

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
