package handlers

import (
	"fmt"
	"net/http"
)

// Does a partial redirect - only reloading the `target`
func HtmxSoftRedirect(w http.ResponseWriter, location string, target string) {
	w.Header().Add("Hx-Location", fmt.Sprintf(`{ "path": "%v", "target": "%v" }`, location, target))
}

// Does a full page redirect
func HtmxHardRedirect(w http.ResponseWriter, location string) {
	w.Header().Add("Hx-Redirect", location)
}
