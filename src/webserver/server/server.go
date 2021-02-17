package server

import (
	"net/http"
	"strings"

	"webserver/server/handlers"

	"github.com/felixge/httpsnoop"
)

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIP
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			url:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)
		m := httpsnoop.CaptureMetrics(h, w, r)
		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration
		logHTTPReq(ri)
	}
	return http.HandlerFunc(fn)
}

// NewHTTPServer creates new http multiplexer
func NewHTTPServer(wd string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RedirectHandler)
	mux.HandleFunc("/cities", handlers.CitiesHandlers)             // shows all user's monitored cities and allows to show weather chart
	mux.HandleFunc("/worker", handlers.WorkersHandler)             // shows all user's workers
	mux.HandleFunc("/worker/new", handlers.AddWorkerHandler)       // adds new worker
	mux.HandleFunc("/worker/delete", handlers.DeleteWorkerHandler) // deletes worker
	mux.HandleFunc("/worker/edit", handlers.EditWorkerHandler)     // edits worker
	mux.HandleFunc("/worker/pause", handlers.PauseWorkerHandler)   // pauses worker
	mux.HandleFunc("/worker/start", handlers.StartWorkerHandler)   // unpause worker
	mux.HandleFunc("/worker/chart", handlers.ChartHandler)         // displays weather chart

	mux.HandleFunc("/login", handlers.LoginHandler)              // login page
	mux.HandleFunc("/logout", handlers.LogoutHandler)            // logout page
	mux.HandleFunc("/signup", handlers.SignupHandler)            // signup page
	mux.HandleFunc("/health-check", handlers.HealthCheckHandler) // signup page

	// redirects resources
	mux.HandleFunc("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(wd+"/ext/"))).ServeHTTP)

	handler := logRequestHandler(mux)
	handler = handlers.SessionManager.LoadAndSave(handler)

	return handler
}
