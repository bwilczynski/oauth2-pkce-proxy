package handlers

import (
	"log"
	"net/http"
)

type accessTokenHandler struct {
	log   *log.Logger
	store ChallengeStore
}

func NewAccessTokenHandler(log *log.Logger, store ChallengeStore) *accessTokenHandler {
	return &accessTokenHandler{log, store}
}

func (h *accessTokenHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.log.Print("AccessToken handler")

	// TODO: validate received code challenge with code verifier
	// h.store.Get(code)

	rw.Write([]byte("access_token"))
}
