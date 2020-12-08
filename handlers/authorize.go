package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type Authorize struct {
	log *log.Logger
}

type AuthorizeRequest struct {
	ClientId            string
	CodeChallenge       string
	CodeChallengeMethod string
	RedirectUri         string
}

const (
	CodeChallengeMethodS256 = "S256"
)

func (ar *AuthorizeRequest) FromQueryParams(r *http.Request) {
	query := r.URL.Query()

	ar.ClientId = query.Get("client_id")
	ar.CodeChallenge = query.Get("code_challenge")
	if codeChallengeMethod := query.Get("code_challenge_method"); codeChallengeMethod != "" {
		ar.CodeChallengeMethod = codeChallengeMethod
	} else {
		ar.CodeChallengeMethod = CodeChallengeMethodS256
	}
	ar.RedirectUri = query.Get("redirect_uri")
}

func NewAuthorize(log *log.Logger) *Authorize {
	return &Authorize{log}
}

func (h *Authorize) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar := &AuthorizeRequest{}
	ar.FromQueryParams(r)

	h.log.Printf("Authorize handler called: %#v", ar)

	if ar.CodeChallengeMethod != CodeChallengeMethodS256 {
		h.log.Printf("Code challenge method %v not supported", ar.CodeChallengeMethod)
		http.Error(w, "Code challenge method not supported", http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("%#v", ar)))
}
