package handlers

import (
	"log"
	"net/http"
)

type accessTokenHandler struct {
	log *log.Logger
}

func NewAccessTokenHandler(log *log.Logger) *accessTokenHandler {
	return &accessTokenHandler{log}
}

func (h *accessTokenHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.log.Print("AccessToken handler")

	// TODO: validate received code challenge with code verifier

	rw.Write([]byte("access_token"))
}
