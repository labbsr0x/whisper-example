package main

import (
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper-client/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func main() {
	// retrieve scopes
	scopesArray := strings.Split(scopes, ",")
	for _, scope := range scopesArray {
		scope = strings.Trim(scope, " ")
	}

	// Initiate the whisper client
	wc := new(client.WhisperClient).
		InitFromParams(whisperURL, clientID, clientSecret, loginRedirectURI, logoutRedirectURI, scopesArray)

	// Register itself in hydra
	token, err := wc.CheckCredentials()
	if err != nil {
		panic("unable to connect w")
	}

	// Add whisper information to context for handlers to use whisper
	ctx := context{
		whisper: whisper{
			client: wc,
			oauthToken: token,
		},
	}

	// Create the necessary routes of the application
	rtr := createRoutes(&ctx)

	// Configure server
	srv := &http.Server{
		Handler: rtr,
		Addr: ":8001",
	}

	// Start server
	err = srv.ListenAndServe()
	gohtypes.PanicIfError("Unable to listen and serve", http.StatusInternalServerError, err)
	logrus.Info("Server started!")
}
