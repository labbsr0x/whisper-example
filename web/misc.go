package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper-client/client"
)

const (
	basePath          = "./assets/html/"
	homePageFile      = "home.html"
	dashboardPageFile = "dashboard.html"
)

// writePage loads a page using templates
func writePage(w http.ResponseWriter, basePath, pageName string, page interface{}) {
	buf := new(bytes.Buffer)
	content := template.Must(template.ParseFiles(basePath + pageName))

	err := content.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load page", http.StatusInternalServerError, err)

	_, err = w.Write(buf.Bytes())
	gohtypes.PanicIfError("Unable to render", http.StatusInternalServerError, err)
}

// getWhisperClient initiate the whisper client
func getWhisperClient() *client.WhisperClient {
	clientID := "client"
	clientSecret := "secret"
	scopes := []string{"openid offline"}
	loginRedirectURI := "http://localhost:8001/dashboard"  // where it should go when finishing authentication
	logoutRedirectURI := "http://localhost:8001/logout" // where it should go when finishing deauthentication
	whisperURL := "http://localhost:7070"                  // whisper path for communication

	return new(client.WhisperClient).InitFromParams(whisperURL, clientID, clientSecret, loginRedirectURI, logoutRedirectURI, scopes)
}

// getWhisperToken retrieve token to be used for authentication inside whisper
func getWhisperToken(whisper *client.WhisperClient) string {
	token, err := whisper.CheckCredentials()
	if err != nil {
		panic("Unable to connect to whisper client")
	}

	tokenString := whisper.GetTokenAsJSONStr(token)

	if tokenString == "" {
		panic("Unable to extract token")
	}

	return tokenString
}

// setHydraCookie set a simple cookie
func setHydraCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})
}

// removeHydraCookie remove a cookie
func removeHydraCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
