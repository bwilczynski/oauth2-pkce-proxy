package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	authorizeURL, _ = url.Parse("http://whatever/authorize")
	tokenURL, _     = url.Parse("http://whatever/token")
)

func TestAuthorizeHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/authorize", nil)
	h := createAuthorizeHandler()

	h.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}
	if loc := rr.Result().Header.Get("Location"); loc == "" {
		t.Errorf("handler returned should return Location header")
	}
}

func TestAuthorizeCodeHandler(t *testing.T) {
	ru := "http://callback"
	rr := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("/code?redirect_uri=%s", ru), nil)
	h := createAuthorizeCodeHandler()

	h.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}
	if loc := rr.Result().Header.Get("Location"); loc != ru {
		t.Errorf("handler returned wrong Location header: got %v want %v", loc, ru)
	}
}

func createAuthorizeHandler() *AuthorizeHandler {
	l := log.New(os.Stdout, "", log.LstdFlags)
	p := &OAuth2Provider{authorizeURL, tokenURL}

	return NewAuthorizeHandler(l, p, "")
}

func createAuthorizeCodeHandler() *AuthorizeCodeHandler {
	l := log.New(os.Stdout, "", log.LstdFlags)

	return NewAuthorizeCodeHandler(l)
}
