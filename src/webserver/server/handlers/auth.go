package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"shared"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/sqlite3store"
	"golang.org/x/crypto/bcrypt"
)

// SessionManager keeps user's sessions
var SessionManager *scs.SessionManager

func init() {
	SessionManager = scs.New()
	SessionManager.Lifetime = 24 * time.Hour
	SessionManager.Store = sqlite3store.New(shared.Db.DB)
	SessionManager.Cookie.Persist = false
}

type authContext struct {
	Title     string
	AuthTitle string
	Action    string
	Submit    string
}

func getLoginContext() authContext {
	return authContext{
		Title:     "Login",
		AuthTitle: "Login",
		Action:    "/login",
		Submit:    "Login",
	}
}

// LoginHandler handles login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/worker", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		var err error
		var t *template.Template

		w.Header().Set("Content-Type", "text/html")
		if t, err = template.ParseFiles(
			"html/auth.html",
			"html/head.html"); err != nil {
			httperr(w, err, http.StatusInternalServerError)
			return
		}
		t.Execute(w, getLoginContext())
		return
	}

	// POST
	var user shared.User
	var err error

	username := r.PostFormValue("username")
	password := []byte(r.PostFormValue("password"))

	// authenticate user
	if user, err = shared.Db.GetUserByName(username); err != nil {
		httperr(w, err, http.StatusUnauthorized)
		return
	}
	if err = bcrypt.CompareHashAndPassword(user.Password, password); err != nil {
		httperr(w, err, http.StatusUnauthorized)
		return
	}

	// create new session
	if err = SessionManager.RenewToken(r.Context()); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	SessionManager.Put(r.Context(), "userid", user.ID)
	SessionManager.RememberMe(r.Context(), true)
	http.Redirect(w, r, "/worker", http.StatusSeeOther)
}

// LogoutHandler handles logout page
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method != "POST" {
		httperr(w, fmt.Errorf("bad HTTP method"), http.StatusBadRequest)
		return
	}

	SessionManager.Remove(r.Context(), "userid")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func getSignupContext() authContext {
	return authContext{
		Title:     "Register",
		AuthTitle: "Register",
		Action:    "/signup",
		Submit:    "Register",
	}
}

// SignupHandler handles signup page
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if SessionManager.Exists(r.Context(), "userid") {
		http.Redirect(w, r, "/worker", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		var err error
		var t *template.Template

		w.Header().Set("Content-Type", "text/html")
		if t, err = template.ParseFiles(
			"html/auth.html",
			"html/head.html"); err != nil {
			httperr(w, err, http.StatusInternalServerError)
			return
		}
		t.Execute(w, getSignupContext())
		return
	}

	// POST
	var err error
	var user shared.User
	user.Username = r.PostFormValue("username")

	if user.Password, err = bcrypt.GenerateFromPassword([]byte(r.PostFormValue("password")), bcrypt.DefaultCost); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}
	if _, err = shared.Db.AddUser(user); err != nil {
		httperr(w, err, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
