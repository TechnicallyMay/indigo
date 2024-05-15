package handlers

import "net/http"

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "settings", []string{"base"}, nil)
}

