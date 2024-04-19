package main

import "net/http"

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "tmpl/products.html")
}

