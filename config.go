package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	ListenPort int
	Provider   *OAuth2Provider
}

type OAuth2Provider struct {
	AuthorizeURL *url.URL
	TokenURL     *url.URL
}

func (cfg *Config) ReadFromEnv(log *log.Logger) {
	if lp := os.Getenv("PKCE_PROXY_PORT"); lp != "" {
		cfg.ListenPort, _ = strconv.Atoi(lp)
	}
	if cfg.ListenPort == 0 {
		cfg.ListenPort = 8080
	}

	authURL, err := url.Parse(os.Getenv("PKCE_PROXY_OAUTH2_AUTH_URL"))
	if err != nil {
		log.Fatal(err)
	}
	tokenURL, err := url.Parse(os.Getenv("PKCE_PROXY_OAUTH2_TOKEN_URL"))
	if err != nil {
		log.Fatal(err)
	}
	cfg.Provider = &OAuth2Provider{authURL, tokenURL}
}
