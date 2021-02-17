package handlers

import (
	"net/http"
)

// RedirectHandler redirects to best fitted page
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/worker", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
