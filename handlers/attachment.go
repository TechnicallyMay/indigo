package handlers

import (
	"io"
	"net/http"
	"os"
)

func AddAttachment(w http.ResponseWriter, path string, name string) {
	w.Header().Set("Content-Disposition", "attachment; filename="+name)

	file, err := os.Open(path)
	handleHttpError(w, err, 500)
	_, err = io.Copy(w, file)
	handleHttpError(w, err, 500)
}
