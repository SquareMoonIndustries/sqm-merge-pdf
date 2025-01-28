package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Route struct for the service
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	UseAuth     bool
}

// Routes for the servcie web handlers
type Routes []Route

// NewRouter creates a new web handler
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = wwwLogger(handler, route.Name)
		if route.UseAuth {
			handler = authHandler(handler, route.Name)
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func authHandler(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			//fmt.Println(username, password)
			hashPassword, user_id := "", 0
			if err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&user_id, &hashPassword); err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				logger.Error(err)
				//fmt.Println("Error: ", err)
				return
			}
			//fmt.Println(user_id)
			_ = user_id
			//fmt.Println(hashPassword)

			// Add subtle.ConstantTimeCompare protection
			err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
			if err == nil {
				inner.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func wwwLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if settings.Debug {
			logger.Info(name + " " + r.RequestURI + " " + r.RemoteAddr + " " + r.Method)
		}
		w.Header().Set("X-Version", appVersionStr)
		inner.ServeHTTP(w, r)
	})
}
