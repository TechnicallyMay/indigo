package handlers

import "net/http"

func GetReportsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, newRenderOpts("reports", nil))
}
