package main

import (
	"net/http"

	"github.com/labbsr0x/goh/gohtypes"
)

// homeHandler renders the home page inserting the url to the whisper login page
func homeHandler(ctx *context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url, err := ctx.whisperClient.GetOAuth2LoginURL()
		gohtypes.PanicIfError("Unable to load redirect url", http.StatusInternalServerError, err)

		pageContent := homePage{LoginURL: url}
		writePage(w, basePath, homePageFile, pageContent)

	}
}

// dashboardHandler renders the dashboard page
// It will try to insert a cookie if there isn't any to authenticate the user. It will panic if it is unable to
// connect to hydra and redirect to home page if it is unable to authorize.
func dashboardHandler(ctx *context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var tokenString string

		if cookie, err := r.Cookie("HAIL_HYDRA"); err != nil {
			code := r.URL.Query().Get("code")

			if code == "" {
				http.Redirect(w, r, "/", http.StatusUnauthorized)
				return
			}

			token, err := ctx.whisperClient.ExchangeCodeForToken(code)
			gohtypes.PanicIfError("Unable to exchange code for token", http.StatusInternalServerError, err)

			tokenString = token.AccessToken

			setHydraCookie(w, tokenString)
		} else {
			tokenString = cookie.Value
		}

		token, err := ctx.whisperClient.IntrospectToken(tokenString)
		gohtypes.PanicIfError("Unable to introspect token", http.StatusInternalServerError, err)

		if !token.Active {
			removeHydraCookie(w)
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}

		url, err := ctx.whisperClient.GetOAuth2LogoutURL()
		gohtypes.PanicIfError("Unable to retrieve logout url", http.StatusInternalServerError, err)

		pageContent := dashboardPage{Username: token.Subject, LogoutURL: url}
		writePage(w, basePath, dashboardPageFile, pageContent)

	}
}
