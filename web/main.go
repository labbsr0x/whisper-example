package main

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labbsr0x/goh/gohtypes"
	"github.com/labbsr0x/whisper-client/client"
)

const (
	indexPage = "index.html"
	dashboardPage = "dashboard.html"
	whisperURL = "http://localhost:7070"
)

var (
	whisperClient *client.WhisperClient
	whisperToken string
)

func writePage(w http.ResponseWriter, pageName string, page interface{}) {
	buf := new(bytes.Buffer)
	content := template.Must(template.ParseFiles("./assets/html/" + pageName))

	err := content.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load page", http.StatusInternalServerError, err)

	_, err = w.Write(buf.Bytes())
	gohtypes.PanicIfError("Unable to render", http.StatusInternalServerError, err)
}

func connectWhisper() (*client.WhisperClient, string) {
	clientID := "client"
	clientSecret := "secret"
	scopes := []string{"openid offline"}
	redirectURI := "http://localhost:8001/login"

	whisper := new(client.WhisperClient).InitFromParams(whisperURL, clientID, clientSecret, redirectURI, scopes)

	token, err := whisper.CheckCredentials()

	if err != nil {
		panic("Unable to connect to whisper client")
	}

	tokenString := whisper.GetTokenAsJSONStr(token)

	if tokenString == "" {
		panic("Unable to extract token")
	}

	return whisper, tokenString
}

func main() {
	whisperClient, _ = connectWhisper()

	rtr := mux.NewRouter()

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		url, err := whisperClient.GetOAuth2LoginURL()

		gohtypes.PanicIfError("Unable to load redirect url", http.StatusInternalServerError, err)

		pageContent := struct {
			Redirect string
		}{
			url,
		}

		writePage(w, indexPage, pageContent)
	}

	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		if cookie, err := r.Cookie("HAIL_HYDRA"); err != nil { // set cookie
			code := r.URL.Query().Get("code")

			if code == "" {
				gohtypes.Panic("Unable to exchange code for token", 6)
			}

			token, err := whisperClient.ExchangeCodeForToken(code)
			gohtypes.PanicIfError("Unable to exchange code for token", 6, err)

			tokenString = token.AccessToken

			http.SetCookie(w, &http.Cookie{
				Name:  "HAIL_HYDRA",
				Value: tokenString,
			})
		} else {
			tokenString = cookie.Value
		}

		token, err := whisperClient.IntrospectToken(tokenString)
		gohtypes.PanicIfError("Unable to introspect token", http.StatusInternalServerError, err)

		if !token.Active {
			http.Redirect(w, nil, "/", http.StatusUnauthorized)
			return
		}

		pageContent := struct {
			Username string
		}{
			token.Subject,
		}

		writePage(w, dashboardPage, pageContent)
	}

	rtr.HandleFunc("/", homeHandler).Methods("GET")
	rtr.HandleFunc("/login", loginHandler).Methods("GET")

	srv := &http.Server{Handler: rtr, Addr: ":8001"}

	err := srv.ListenAndServe()
	gohtypes.PanicIfError("Unable to listen and serve", http.StatusInternalServerError, err)
}
