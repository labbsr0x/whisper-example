package main

import (
	"github.com/labbsr0x/whisper-client/client"
	"net/http"

	"github.com/labbsr0x/goh/gohtypes"
)

// homeHandler renders the home page inserting the url to the whisper login page
func homeHandler(ctx *context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if cookie, err := r.Cookie("ACCESS_TOKEN"); err == nil {
			token, err := ctx.whisperClient.IntrospectToken(cookie.Value)
			gohtypes.PanicIfError("Unable to introspect token", http.StatusInternalServerError, err)

			if token.Active {
				http.Redirect(w, r, "/dashboard", http.StatusFound)
				return
			}
		}

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

		var tokens client.Tokens

		if cookie, err := r.Cookie("ACCESS_TOKEN"); err != nil {
			code := r.URL.Query().Get("code")

			if code == "" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			tokens, err = ctx.whisperClient.ExchangeCodeForToken(code)
			gohtypes.PanicIfError("Unable to exchange code for token", http.StatusInternalServerError, err)

			setHydraCookie(w, "ACCESS_TOKEN", tokens.AccessToken)
			setHydraCookie(w, "OPENID_TOKEN", tokens.OpenIdToken)
			setHydraCookie(w, "REFRESH_TOKEN", tokens.RefreshToken)

			http.Redirect(w, r, "/dashboard", http.StatusFound) // redirect to self without the parameters
			return
		} else {
			tokens.AccessToken = cookie.Value

			cookie, _ = r.Cookie("OPENID_TOKEN")
			tokens.OpenIdToken = cookie.Value

			cookie, _ = r.Cookie("REFRESH_TOKEN")
			tokens.RefreshToken = cookie.Value
		}

		token, err := ctx.whisperClient.IntrospectToken(tokens.AccessToken)
		gohtypes.PanicIfError("Unable to introspect token", http.StatusInternalServerError, err)

		if !token.Active {
			logoutHandler(w, r)
			return
		}

		url, err := ctx.whisperClient.GetOAuth2LogoutURL(tokens.OpenIdToken, "http://localhost:8001/logout")
		gohtypes.PanicIfError("Unable to retrieve logout url", http.StatusInternalServerError, err)

		pageContent := dashboardPage{Username: token.Subject, LogoutURL: url}
		writePage(w, basePath, dashboardPageFile, pageContent)
	}
}

// logoutHandler remove cookies and redirect to logout
func logoutHandler (w http.ResponseWriter, r *http.Request) {
	removeHydraCookie(w, "ACCESS_TOKEN")
	removeHydraCookie(w, "OPENID_TOKEN")
	removeHydraCookie(w, "REFRESH_TOKEN")

	http.Redirect(w, r, "/", http.StatusFound)
}