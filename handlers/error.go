package handlers

import (
	"log"
	"net/http"
)

func handleHttpError(w http.ResponseWriter, err error, code int) {
	if err == nil {
		return
	}
	log.Println("ERROR:", err)
	http.Error(w, err.Error(), code)
}
