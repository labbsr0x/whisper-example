package main

import (
	"github.com/labbsr0x/goh/gohtypes"
	"net/http"
)

// postLoginHandler set the necessary cookies to identify the user. It will exchange the code in the request for tokens,
// set them as cookies and redirect to the specified url.
func postLoginHandler(ctx *context, redirectTo string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve the code query param
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// exchange the code retrieved by tokens that identify the user and more
		codeVerifierCookie, err := r.Cookie("CODE_VERIFIER")
		if err != nil {
			return
		}
		stateCookie, err := r.Cookie("STATE")
		if err != nil {
			return
		}

		tokens, err := ctx.whisper.client.ExchangeCodeForToken(code, codeVerifierCookie.Value, stateCookie.Value)
		gohtypes.PanicIfError("Unable to exchange code for token", http.StatusInternalServerError, err)

		// set all information related to user identity in cookies
		setCookie(w, accessTokenCookieName, tokens.AccessToken)
		setCookie(w, openIDTokenCookieName, tokens.OpenIdToken)
		setCookie(w, refreshTokenCookieName, tokens.RefreshToken)

		// Redirect to the specified url
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}

// postLogoutHandler overwrite cookies used to identify the user and redirect to the specified url
func postLogoutHandler(redirectTo string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// all cookies related to the user identity should be removed. Hydra does not invalidate the tokens
		// stored inside the cookies, it just forget sessions.
		unsetCookie(w, accessTokenCookieName)
		unsetCookie(w, openIDTokenCookieName)
		unsetCookie(w, refreshTokenCookieName)

		// Redirect to the specified url
		http.Redirect(w, r, redirectTo, http.StatusFound)
	}
}
