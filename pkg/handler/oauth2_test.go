package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestOAuth2LoginInvalidProvider(t *testing.T) {

	req, err := http.NewRequest("GET", "/login/invalidprovider", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPOAuth2(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOAuth2LoginNoConfig(t *testing.T) {

	providers["google"].ClientID = ""
	providers["google"].ClientSecret = ""
	providers["google"].RedirectURL = ""

	req, err := http.NewRequest("GET", "/login/google", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPOAuth2(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOAuth2LoginWithConfig(t *testing.T) {

	providers["google"].ClientID = "client_id"
	providers["google"].ClientSecret = "secret_id"
	providers["google"].RedirectURL = "redirect_url"

	req, err := http.NewRequest("GET", "/login/google", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPOAuth2(req)

	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
}

func TestOAuth2CallbackInvalidProvider(t *testing.T) {

	req, err := http.NewRequest("POST", "/callback/invalidprovider", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPOAuth2(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOAuth2CallbackNoConfig(t *testing.T) {

	providers["google"].ClientID = ""
	providers["google"].ClientSecret = ""
	providers["google"].RedirectURL = ""

	req, err := http.NewRequest("POST", "/callback/google", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPOAuth2(req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOAuth2CallbackWithConfig(t *testing.T) {

	providers["google"].ClientID = "client_id"
	providers["google"].ClientSecret = "secret_id"
	providers["google"].RedirectURL = "redirect_url"

	req, err := http.NewRequest("POST", "/callback/google", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPOAuth2(req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func serveHTTPOAuth2(req *http.Request) *httptest.ResponseRecorder {

	rec := httptest.NewRecorder()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login/{provider:[a-z]+}", OAuth2Login).Methods("GET")
	router.HandleFunc("/callback/{provider:[a-z]+}", OAuth2Callback).Methods("POST")
	router.ServeHTTP(rec, req)

	return rec
}
