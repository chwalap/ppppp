package handlers

import (
	"fmt"
	"net/http"

	"shared"
)

// Cities contains all available cities
var Cities map[int64]shared.City

func init() {
	var err error
	Cities, err = shared.LoadCities()
	if err != nil {
		panic(err)
	}
}

// CitiesHandlers handles displaing cities
func CitiesHandlers(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		httperr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}
}
