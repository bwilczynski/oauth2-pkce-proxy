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

func writeError(w http.ResponseWriter, res *m.ValidationResult) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
}
