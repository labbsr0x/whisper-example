package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labbsr0x/goh/gohtypes"
)

func main() {
	whisperClient := getWhisperClient()
	whisperToken := getWhisperToken(whisperClient)

	ctx := context{
		whisperClient: whisperClient,
		whisperToken:  whisperToken,
	}

	rtr := mux.NewRouter()

	rtr.HandleFunc("/", homeHandler(&ctx)).
		Methods("GET")

	rtr.HandleFunc("/dashboard", dashboardHandler(&ctx)).
		Methods("GET")

	srv := &http.Server{Handler: rtr, Addr: ":8001"}

	err := srv.ListenAndServe()
	gohtypes.PanicIfError("Unable to listen and serve", http.StatusInternalServerError, err)
}
