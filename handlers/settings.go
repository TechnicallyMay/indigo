package handlers

import "net/http"

func HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, newRenderOpts("settings", nil))
}
