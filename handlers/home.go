package handlers

import "net/http"

func GetRootHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/home/", 308)
}

func GetHomeHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "home", []string{"base"})
}
