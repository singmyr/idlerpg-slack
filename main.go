package main

import (
	"fmt"
	"net/http"
)

func authFunc(w http.ResponseWriter, r *http.Request) {
	//authorizeUrl := "https://slack.com/oauth/authorize"
	//clientID := ""
	// Scopes: https://api.slack.com/docs/oauth-scopes
	//scope := "bot"
	//redirectUri := ""
	//state := ""
	//team := ""
	fmt.Fprintf(w, "<h1>%s</h1>", "Hello world!")
}

//func main() {
//	http.HandleFunc("/", authFunc)
//	fmt.Println("HTTP Server started")
//	log.Fatal(http.ListenAndServe(":80", nil))
//}
