package handlers

import (
	"net/http"
)

// HealthCheckHandler is available when this server is
func HealthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
