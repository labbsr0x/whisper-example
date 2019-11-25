package main

const (
	clientID         = "client"
	clientSecret     = "secret"
	scopes           = "openid,offline" // separated by commas
	whisperURL       = "http://localhost:7070" // whisper path for communication
	postLoginPath = "/post-login"
	postLogoutPath = "/post-logout"
	host             = "http://localhost:8001"
	loginRedirectURI = host + postLoginPath
	logoutRedirectURI  = host + postLogoutPath

	accessTokenCookieName  = "WHISPER_ACCESS_TOKEN"
	openIDTokenCookieName  = "WHISPER_OPENID_TOKEN"
	refreshTokenCookieName = "WHISPER_REFRESH_TOKEN"
)