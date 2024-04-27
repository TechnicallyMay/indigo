package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/TechnicallyMay/indigo/handlers"
)
var validPath = regexp.MustCompile("^(/[a-zA-Z0-9]+/{0,1}?)*$")

func main() {
    mux := http.NewServeMux()
    mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

    mux.HandleFunc("GET /{$}", handlers.GetRootHandler)
    mux.HandleFunc("GET /home/", handlers.GetHomeHandler)
    mux.HandleFunc("GET /billing/", handlers.GetBillingHandler)
    mux.HandleFunc("GET /customers/", handlers.GetCustomersHandler)
    mux.HandleFunc("GET /products/", handlers.GetProductsHandler)
    mux.HandleFunc("GET /settings/", handlers.GetSettingsHandler)

    server := &http.Server{
        Addr: ":8080",
        Handler: mux,
    }

    log.Fatal(server.ListenAndServe())
}

