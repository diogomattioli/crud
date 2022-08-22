package handler

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginOk(t *testing.T) {

	SetAuthenticator(&MockAuth{})

	body, header, err := formData(map[string]string{"user": "a", "pass": "a"})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", header)

	rec := serveHTTPAuth(req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "123-token", rec.Header().Get("X-Access-Token"))
}

func TestLoginNoLoginPass(t *testing.T) {

	SetAuthenticator(&MockAuth{})

	body, header, err := formData(map[string]string{"user": "a"})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", header)

	rec := serveHTTPAuth(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	body, header, err = formData(map[string]string{"pass": "a"})
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", "/login/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", header)

	rec = serveHTTPAuth(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req, err = http.NewRequest("POST", "/login/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec = serveHTTPAuth(req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginAuthenticaionFail(t *testing.T) {

	SetAuthenticator(&MockAuth{shouldFail: true})

	body, header, err := formData(map[string]string{"user": "a", "pass": "a"})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", header)

	rec := serveHTTPAuth(req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthListOk(t *testing.T) {

	SetAuthenticator(&MockAuth{})

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/auth/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("X-Access-Token", "123-token")

	rec := serveHTTPAuth(req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthWrongToken(t *testing.T) {

	SetAuthenticator(&MockAuth{})

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/auth/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("X-Access-Token", "token")

	rec := serveHTTPAuth(req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAuthEmpty(t *testing.T) {

	SetAuthenticator(&MockAuth{})

	setupDb(1)
	defer destroyDb()

	req, err := http.NewRequest("GET", "/auth/dummy/", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rec := serveHTTPAuth(req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
