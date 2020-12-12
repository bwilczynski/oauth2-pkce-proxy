package handlers

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

var provider *OAuth2Provider

func init() {
	authURL, err := url.Parse(os.Getenv("PKCE_PROXY_OAUTH2_AUTH_URL"))
	if err != nil {
		log.Fatal(err)
	}

	tokenURL, err := url.Parse(os.Getenv("PKCE_PROXY_OAUTH2_TOKEN_URL"))
	if err != nil {
		log.Fatal(err)
	}

	provider = NewOAuth2Provider(authURL, tokenURL)
}

func RegisterRoutes(mux *http.ServeMux, log *log.Logger) {
	mux.Handle("/authorize", NewAuthorize(log, provider, "/code"))
	mux.Handle("/access_token", NewAccessToken(log))
	mux.Handle("/code", NewAuthorizeCode(log))
}
