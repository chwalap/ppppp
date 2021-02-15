package handlers

import (
	"log"
	"net/http"
)

func httperr(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
	log.Println("Error:", err.Error(), "->", http.StatusText(status))
}
