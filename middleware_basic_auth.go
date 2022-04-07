package az

import (
	"fmt"
	"log"
	"net/http"
)

// Param 1: username; Param 2: password
func basicAuth(h http.HandlerFunc, params ...interface{}) http.HandlerFunc {
	if len(params) != 2 {
		log.Fatal("Param 1: username; Param 2: password not provided")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == fmt.Sprint(params[0]) && password == fmt.Sprint(params[1]) {
			// Delegate request to the given handle
			h(w, r)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}
