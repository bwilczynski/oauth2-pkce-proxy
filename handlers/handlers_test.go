package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	ar := &models.AuthorizeRequest{CodeChallenge: "whatever", ClientId: "whatever", RedirectUri: "whatever"}
	r, _ := http.NewRequest("GET", fmt.Sprintf("/authorize?%s", ar.URLQuery()), nil)
	h := createAuthorizeHandler()

	h.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.NotEqual(t, "", rr.Result().Header.Get("Location"))
}

func createAuthorizeHandler() *authorizeHandler {
	l := zerolog.Nop()
	p := &models.OAuth2Provider{AuthorizationEndpoint: "http://whatever/authorize", TokenEndpoint: "http://whatever/token"}

	return NewAuthorizeHandler(&l, p, "")
}

func TestAuthorizeCodeHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	cr := &models.AuthorizeCodeRequest{Code: "whatever", RedirectUri: "http://callback"}
	r, _ := http.NewRequest("GET", fmt.Sprintf("/code?%s", cr.URLQuery()), nil)
	h := createAuthorizeCodeHandler()

	h.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.True(t, strings.Index(rr.Result().Header.Get("Location"), cr.RedirectUri) == 0)
}

func createAuthorizeCodeHandler() *authorizeCodeHandler {
	l := zerolog.Nop()
	store := NewInMemoryChallengeStore()

	return NewAuthorizeCodeHandler(&l, store)
}
