package handlers

import "net/http"

func GetBillingHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "billing", []string{"base"}, nil)
}
