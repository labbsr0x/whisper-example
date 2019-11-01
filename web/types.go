package main

import "github.com/labbsr0x/whisper-client/client"

// context holds the context of the application
type context struct {
	whisperClient *client.WhisperClient
	whisperToken  string
}

// homePage holds the info necessary for the home page template
type homePage struct {
	LoginURL string
}

// dashboardPage holds the info necessary for the dashboard page template
type dashboardPage struct {
	Username string
	LogoutURL string
}
