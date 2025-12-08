package handlers

import "net/http"

func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home/", http.StatusMovedPermanently)
}

func GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, newRenderOpts("home", nil))
}
