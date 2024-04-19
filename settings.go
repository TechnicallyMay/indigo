package main

import "net/http"

func getSettingsHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "tmpl/settings.html")
}

