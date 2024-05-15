package handlers

import (
	"fmt"
	"net/http"
)

func HtmxRedirect(w http.ResponseWriter, location string, target string) {
    w.Header().Add("Hx-Location", fmt.Sprintf("{ \"path\": \"%v\", \"target\": \"%v\" }", location, target))
}

