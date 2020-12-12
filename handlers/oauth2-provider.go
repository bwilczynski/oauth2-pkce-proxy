package handlers

import "net/url"

type OAuth2Provider struct {
	AuthorizeURL *url.URL
	TokenURL     *url.URL
}

func NewOAuth2Provider(authorizeURL *url.URL, tokenURL *url.URL) *OAuth2Provider {
	return &OAuth2Provider{authorizeURL, tokenURL}
}
