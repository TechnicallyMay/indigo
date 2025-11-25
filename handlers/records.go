package handlers

import "net/http"

func GetRecordsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "records", []string{"base"}, nil)
}
