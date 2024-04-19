package main

import "net/http"

func getCustomersHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "tmpl/customers.html")
}

