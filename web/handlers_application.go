package main

import (
	"fmt"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper-client/client"
	"net/http"
)


// homeHandler create a function to mount the home page. It should be available only to unidentified users, that is,
// users that are not logged in. If there is a logged in user, it will redirect to the dashboard, the default home for logged in users
// and if there isn't a user logged in, it will retrieve a login url to be used in the page and mount the page
func homeHandler(ctx *context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify if there is a user logged. If there is, it should be able to retrieve valid tokens from the request
		_, _, err := GetTokensFromRequest(w, r, ctx.whisper.client)
		if err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}

		// Mount the home page with the login url that connects with whisper
		url, err := ctx.whisper.client.GetOAuth2LoginURL()
		gohtypes.PanicIfError("Unable to load redirect url", http.StatusInternalServerError, err)

		pageContent := homePage{
			LoginURL: url,
		}

		writePage(w, "home", pageContent)
	}
}

// dashboardHandler create a function to mount the dashboard page. It should be available only to identified users, that
// is, logged in users. If you are not logged, it will redirect to home page, the default home for unidentified users. It
// will panic if it is unable to connect to hydra.
func dashboardHandler(ctx *context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the cookies of the logged user. If there is an error, the current user is not logged
		tokens, token, err := GetTokensFromRequest(w, r, ctx.whisper.client)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Mount the dashboard page with the logout url that connects with whisper
		url, err := ctx.whisper.client.GetOAuth2LogoutURL(tokens.OpenIdToken, logoutRedirectURI)
		gohtypes.PanicIfError("Unable to retrieve logout url", http.StatusInternalServerError, err)

		pageContent := dashboardPage{
			Username: token.Subject,
			LogoutURL: url,
		}

		writePage(w, "dashboard", pageContent)
	}
}

// GetTokensFromRequest extract tokens from request and introspect the access token with hydra. From the request it will
// identify the cookies set on post login and extract them. If they are valid, it means that there is a logged user and
// the lone token returned contains its username.
func GetTokensFromRequest(w http.ResponseWriter, r *http.Request, wc *client.WhisperClient) (tokens client.Tokens, token client.Token, err error) {
	// retrieve access token
	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return
	}
	tokens.AccessToken = cookie.Value

	// retrieve open id token
	cookie, err = r.Cookie(openIDTokenCookieName)
	if err != nil {
		return
	}
	tokens.OpenIdToken = cookie.Value

	// retrieve refresh token
	cookie, err = r.Cookie(refreshTokenCookieName)
	if err != nil {
		return
	}
	tokens.RefreshToken = cookie.Value

	// Verify if the access token in valid
	if token, err = wc.IntrospectToken(tokens.AccessToken); err != nil || !token.Active {
		if err == nil {
			err = fmt.Errorf("invalid token")
		}

		unsetCookie(w, accessTokenCookieName)
		unsetCookie(w, openIDTokenCookieName)
		unsetCookie(w, refreshTokenCookieName)
	}

	return
}

