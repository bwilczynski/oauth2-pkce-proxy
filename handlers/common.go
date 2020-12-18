package handlers

import (
	"encoding/json"
	"net/http"

	m "github.com/bwilczynski/oauth2-pkce-proxy/models"
)

const (
	challengeCookieName   = "pkce-proxy-challenge"
	challengeCookieMaxAge = 300
)

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	if verr, ok := err.(*m.ValidationError); ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(verr)
	}
}

type ChallengeStore interface {
	Write(code string, verifier m.CodeVerifier)
	Get(code string) m.CodeVerifier
}

type inMemoryChallengeStore struct {
	values map[string]m.CodeVerifier
}

func NewInMemoryChallengeStore() ChallengeStore {
	return &inMemoryChallengeStore{values: make(map[string]m.CodeVerifier)}
}

func (ms *inMemoryChallengeStore) Write(code string, verifier m.CodeVerifier) {
	ms.values[code] = verifier
}

func (ms *inMemoryChallengeStore) Get(code string) m.CodeVerifier {
	return ms.values[code]
}
