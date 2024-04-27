package handlers

import "net/http"

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "customers", []string{"base"})
}

