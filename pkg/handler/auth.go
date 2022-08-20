package handler

import (
	"net/http"

	"github.com/diogomattioli/crud/pkg/data"
)

var auth data.Authenticator

func SetAuthenticator(_auth data.Authenticator) {
	auth = _auth
}

func Login(w http.ResponseWriter, r *http.Request) {

	user := r.FormValue("user")
	pass := r.FormValue("pass")

	if user == "" || pass == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !auth.Authenticate(user, pass) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := auth.Create()

	w.Header().Set("X-Access-Token", token)
}

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("X-Access-Token")

		if !auth.Use(token) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
