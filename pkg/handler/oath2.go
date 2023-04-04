package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

var (
	providers = map[string]*oauth2.Config{
		"google": {
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		},
		"microsoft": {
			ClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
			ClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("MICROSOFT_REDIRECT_URL"),
			Endpoint:     microsoft.AzureADEndpoint("common"),
			Scopes:       []string{"user.read"},
		},
		"facebook": {
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
			Endpoint:     facebook.Endpoint,
			Scopes:       []string{"email", "public_profile"},
		},
	}
)

func OAuth2Login(w http.ResponseWriter, r *http.Request) {

	provider := mux.Vars(r)["provider"]

	if providerConf, exists := providers[provider]; exists && providerConf.ClientID != "" && providerConf.ClientSecret != "" && providerConf.RedirectURL != "" {
		url := providerConf.AuthCodeURL("state")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// After the user grants permission to our application, the provider will redirect them back to our callback URL
func OAuth2Callback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")

	// Exchange the authorization code for an access token
	var token *oauth2.Token
	var err error

	provider := mux.Vars(r)["provider"]

	if providerConf, exists := providers[provider]; exists && providerConf.ClientID != "" && providerConf.ClientSecret != "" && providerConf.RedirectURL != "" {
		token, err = providerConf.Exchange(context.Background(), code)
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Use the token to make API requests on behalf of the user
	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
	fmt.Fprintf(w, "Refresh Token: %s\n", token.RefreshToken)
	fmt.Fprintf(w, "Expiry: %s\n", token.Expiry.String())
}
