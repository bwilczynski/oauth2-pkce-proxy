package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bwilczynski/oauth2-pkce-proxy/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeHandler(t *testing.T) {
	testcases := []struct {
		query map[string][]string
		code  int
	}{
		{
			query: map[string][]string{},
			code:  http.StatusBadRequest,
		},
		{
			query: map[string][]string{"redirect_uri": {"whatever"}, "code_challenge": {"whatever"}},
			code:  http.StatusFound,
		},
	}

	for _, tc := range testcases {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", fmt.Sprintf("/authorize?%s", url.Values(tc.query).Encode()), nil)
		h := createAuthorizeHandler()

		h.ServeHTTP(rr, r)

		assert.Equal(t, tc.code, rr.Code)
	}
}

func createAuthorizeHandler() *authorizeHandler {
	l := zerolog.Nop()
	p := &models.OAuth2Provider{AuthorizationEndpoint: "http://whatever/authorize", TokenEndpoint: "http://whatever/token"}

	return NewAuthorizeHandler(&l, p, "")
}

func TestAuthorizeCodeHandler(t *testing.T) {
	testcases := []struct {
		query map[string][]string
		code  int
	}{
		{
			query: map[string][]string{},
			code:  http.StatusBadRequest,
		},
		{
			query: map[string][]string{"code": {"whatever"}, "redirect_uri": {"whatever"}},
			code:  http.StatusFound,
		},
	}

	for _, tc := range testcases {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", fmt.Sprintf("/code?%s", url.Values(tc.query).Encode()), nil)
		h := createAuthorizeCodeHandler()

		h.ServeHTTP(rr, r)

		assert.Equal(t, tc.code, rr.Code)
	}
}

func createAuthorizeCodeHandler() *authorizeCodeHandler {
	l := zerolog.Nop()
	store := NewInMemoryChallengeStore()

	return NewAuthorizeCodeHandler(&l, store)
}
