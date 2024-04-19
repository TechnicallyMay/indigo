package main

import "net/http"

func getInvoicesHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "tmpl/invoices.html")
}

