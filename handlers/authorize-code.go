package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type AuthorizeCode struct {
	log *log.Logger
}

func NewAuthorizeCode(log *log.Logger) *AuthorizeCode {
	return &AuthorizeCode{log}
}

type authorizeCodeRequest struct {
	Code        string
	RedirectUri string
}

func (cr *authorizeCodeRequest) FromQueryParams(r *http.Request) {
	query := r.URL.Query()

	cr.Code = query.Get("code")
	cr.RedirectUri = query.Get("redirect_uri")
}

func (h *AuthorizeCode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cr := &authorizeCodeRequest{}
	cr.FromQueryParams(r)

	h.log.Printf("AuthorizeCode handler called: %#v", cr)

	q := r.URL.Query()
	q.Del("redirect_uri")

	redirectURI := fmt.Sprintf("%s?%s", cr.RedirectUri, q.Encode())

	w.Header().Add("Location", redirectURI)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
