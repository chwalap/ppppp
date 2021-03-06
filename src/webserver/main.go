package main

import (
	"log"
	"net/http"
	"os"

	"webserver/server"
)

func main() {
	var err error
	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		panic(err)
	}

	log.Printf("Webserver service starting { cwd: %s }\n", cwd)
	server.OpenHTTPLog()

	handler := server.NewHTTPServer(cwd)
	http.ListenAndServe(":8080", handler)
}
