package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	challengeCookieName   = "pkce-proxy-challenge"
	challengeCookieMaxAge = 300
)

type AuthorizeHandler struct {
	log          *log.Logger
	provider     *OAuth2Provider
	callbackPath string
}

func NewAuthorizeHandler(log *log.Logger, provider *OAuth2Provider, callbackPath string) *AuthorizeHandler {
	return &AuthorizeHandler{log, provider, callbackPath}
}

func (h *AuthorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ar := &AuthorizeRequest{}
	ar.FromQueryParams(r)

	h.log.Printf("Authorize handler called: %#v", ar)

	if ar.CodeChallengeMethod != CodeChallengeMethodS256 {
		h.log.Printf("Code challenge method %v not supported", ar.CodeChallengeMethod)
		http.Error(w, "Code challenge method not supported", http.StatusBadRequest)
		return
	}

	q := r.URL.Query()
	q.Set("redirect_uri", fmt.Sprintf("http://%s%s?redirect_uri=%s", r.Host, h.callbackPath, ar.RedirectUri))

	redirectURI := fmt.Sprintf("%s?%s", h.provider.AuthorizeURL, q.Encode())
	w.Header().Add("Location", redirectURI)
	setChallengeCookie(w, ar.CodeChallenge)
	w.WriteHeader(http.StatusFound)
}

func setChallengeCookie(w http.ResponseWriter, challenge string) {
	cookie := http.Cookie{Name: challengeCookieName, Value: url.QueryEscape(challenge), Path: "/", HttpOnly: true, MaxAge: challengeCookieMaxAge}
	http.SetCookie(w, &cookie)
}

type AuthorizeCodeHandler struct {
	log *log.Logger
}

func NewAuthorizeCodeHandler(log *log.Logger) *AuthorizeCodeHandler {
	return &AuthorizeCodeHandler{log}
}

func (h *AuthorizeCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cr := &AuthorizeCodeRequest{}
	cr.FromQueryParams(r)
	cc := readChallengeCookie(r)

	h.log.Printf("AuthorizeCode handler called: %#v, cc: %s", cr, cc)

	q := r.URL.Query()
	q.Del("redirect_uri")

	redirectURI := cr.RedirectUri
	if len(q) > 0 {
		redirectURI += "&" + q.Encode()
	}

	w.Header().Add("Location", redirectURI)
	w.WriteHeader(http.StatusFound)
}

func readChallengeCookie(r *http.Request) string {
	c, err := r.Cookie(challengeCookieName)
	if err != nil {
		return ""
	}
	challenge, _ := url.QueryUnescape(c.Value)
	return challenge
}

type AccessTokenHandler struct {
	log *log.Logger
}

func NewAccessTokenHandler(log *log.Logger) *AccessTokenHandler {
	return &AccessTokenHandler{log}
}

func (h *AccessTokenHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.log.Print("AccessToken handler")
	rw.Write([]byte("access_token"))
}
