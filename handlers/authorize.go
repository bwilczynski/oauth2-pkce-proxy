package handlers

import (
	"log"
	"net/http"
)

type Authorize struct {
	log *log.Logger
}

func NewAuthorize(log *log.Logger) *Authorize {
	return &Authorize{log}
}

func (h *Authorize) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.log.Print("Authorize handler")
	rw.Write([]byte("authorize"))
}
