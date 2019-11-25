package main

import (
	"github.com/gorilla/mux"
)

// createRoutes set the routes the application need. That is, it will create the suggested web hooks for whisper and the
// routes of the main application. The suggest web hooks for whisper are to separate some recommended operation from the
// logic of your application.
func createRoutes(ctx *context) *mux.Router {
	rtr := mux.NewRouter()

	// Whisper necessary web hooks
	rtr.HandleFunc(postLoginPath, postLoginHandler(ctx, "/dashboard")).
		Methods("GET")

	rtr.HandleFunc(postLogoutPath, postLogoutHandler("/")).
		Methods("GET")

	// Application routes
	rtr.HandleFunc("/", homeHandler(ctx)).
		Methods("GET")

	rtr.HandleFunc("/dashboard", dashboardHandler(ctx)).
		Methods("GET")

	return rtr
}