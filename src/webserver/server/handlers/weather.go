package handlers

import (
	"fmt"
	"net/http"
)

// WeatherChartHandlers handles weather chart
func WeatherChartHandlers(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		httperr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}
}
