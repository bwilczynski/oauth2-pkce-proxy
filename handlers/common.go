package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

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
	sync.Mutex
	values map[string]m.CodeVerifier
}

func NewInMemoryChallengeStore() *inMemoryChallengeStore {
	return &inMemoryChallengeStore{values: make(map[string]m.CodeVerifier)}
}

func (ms *inMemoryChallengeStore) Write(code string, verifier m.CodeVerifier) {
	ms.Lock()
	ms.values[code] = verifier
	ms.Unlock()
}

func (ms *inMemoryChallengeStore) Get(code string) m.CodeVerifier {
	ms.Lock()
	res := ms.values[code]
	ms.Unlock()
	return res
}
