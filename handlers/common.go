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
	Write(code string, challenge string)
	Get(code string) string
}

type inMemoryChallengeStore struct {
	sync.Mutex
	values map[string]string
}

func NewInMemoryChallengeStore() *inMemoryChallengeStore {
	return &inMemoryChallengeStore{values: make(map[string]string)}
}

func (ms *inMemoryChallengeStore) Write(code string, challenge string) {
	ms.Lock()
	ms.values[code] = challenge
	ms.Unlock()
}

func (ms *inMemoryChallengeStore) Get(code string) string {
	ms.Lock()
	res := ms.values[code]
	ms.Unlock()
	return res
}
