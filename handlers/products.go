package handlers

import "net/http"

func HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "products", nil, nil)
}
