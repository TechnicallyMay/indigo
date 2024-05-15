package handlers

import "net/http"

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "products", []string{"base"}, nil)
}

