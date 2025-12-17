package handlers

import "net/http"

func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/billing", http.StatusMovedPermanently)
}
