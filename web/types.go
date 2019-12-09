package main

import (
	"github.com/labbsr0x/whisper-client/client"
	"golang.org/x/oauth2"
)

// whisper holds the whisper connection and token
type whisper struct {
	client *client.WhisperClient
	oauthToken *oauth2.Token
}

// context holds the context of the application
type context struct {
	whisper whisper
}

// homePage holds the info necessary for the home page template
type homePage struct {
	LoginURL string
}

// dashboardPage holds the info necessary for the dashboard page template
type dashboardPage struct {
	Username  string
	LogoutURL string
}
