package main

import (
	"flag"
	"fmt"
	"github.com/schwiet/slack-spotify/lambda/spotify"
	"net/http"
)

func main() {

	port := flag.Int(
		"port", 8008, "Set port at which to listen for HTTP connections")

	// http Server Mux
	mux := http.NewServeMux()

	// create a cancellable context and cancel when the main loop exits
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	mux.HandleFunc("/auth/", spotify.AuthorizeHandler)

	// start server, this will block until an error is encountered
	addrStr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(addrStr, mux)
	if err != nil {
		fmt.Println("Web Server Exiting with error:", err.Error())
	}
}
