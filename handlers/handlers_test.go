package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/bwilczynski/oauth2-pkce-proxy/models"
)

var (
	authorizeURL, _ = url.Parse("http://whatever/authorize")
	tokenURL, _     = url.Parse("http://whatever/token")
)

func TestAuthorizeHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	ar := &models.AuthorizeRequest{CodeChallenge: "whatever", ClientId: "whatever", RedirectUri: "whatever"}
	r, _ := http.NewRequest("GET", fmt.Sprintf("/authorize?%s", ar.URLQuery()), nil)
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
	rr := httptest.NewRecorder()
	cr := &models.AuthorizeCodeRequest{Code: "whatever", RedirectUri: "http://callback"}
	r, _ := http.NewRequest("GET", fmt.Sprintf("/code?%s", cr.URLQuery()), nil)
	h := createAuthorizeCodeHandler()

	h.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusFound)
	}
	if loc := rr.Result().Header.Get("Location"); strings.Index(loc, cr.RedirectUri) != 0 {
		t.Errorf("handler returned Location header not containing %v: got %v", cr.RedirectUri, loc)
	}
}

func createAuthorizeHandler() *authorizeHandler {
	l := log.New(os.Stdout, "", log.LstdFlags)
	p := &models.OAuth2Provider{AuthorizeURL: authorizeURL, TokenURL: tokenURL}

	return NewAuthorizeHandler(l, p, "")
}

func createAuthorizeCodeHandler() *authorizeCodeHandler {
	l := log.New(os.Stdout, "", log.LstdFlags)
	store := NewInMemoryChallengeStore()

	return NewAuthorizeCodeHandler(l, store)
}
