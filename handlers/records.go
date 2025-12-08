package handlers

import "net/http"

func HandleGetRecords(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, newRenderOpts("records", nil))
}
