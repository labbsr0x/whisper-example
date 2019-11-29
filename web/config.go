package main

const (
	clientID         = "client"
	clientSecret     = "secret"
	scopes           = "openid,offline" // separated by commas

	postLoginPath = "/post-login"
	postLogoutPath = "/post-logout"

	accessTokenCookieName  = "WHISPER_ACCESS_TOKEN"
	openIDTokenCookieName  = "WHISPER_OPENID_TOKEN"
	refreshTokenCookieName = "WHISPER_REFRESH_TOKEN"
)

var (
	whisperURL = "http://localhost:7070" // whisper path for communication
	selfURL    = "http://localhost:8001"

	loginRedirectURI = selfURL + postLoginPath
	logoutRedirectURI  = selfURL + postLogoutPath
)
