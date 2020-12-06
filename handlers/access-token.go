package handlers

import (
	"log"
	"net/http"
)

type AccessToken struct {
	log *log.Logger
}

func NewAccessToken(log *log.Logger) *AccessToken {
	return &AccessToken{log}
}

func (h *AccessToken) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.log.Print("AccessToken handler")
	rw.Write([]byte("access_token"))
}
