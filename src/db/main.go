package main

import (
	"fmt"
	"log"
	"net/http"

	"shared"
)

func main() {
	log.Printf("Database hypervisor started\n")

	if err := shared.Db.CreateDbIfNotExists(false); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", healthCheck)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", shared.DbPort), mux))
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
