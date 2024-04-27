package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/TechnicallyMay/indigo/handlers"
)
var validPath = regexp.MustCompile("^(/[a-zA-Z0-9]+/{0,1}?)*$")

// TODO: Utilize Must and only parse templates once
// var templates = template.Must(template.ParseFiles("tmpl/base.html", "tmpl/view.html", "tmpl/edit.html"))

func main() {
    mux := http.NewServeMux()
    mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    mux.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

    mux.HandleFunc("GET /{$}", handlers.GetRootHandler)
    mux.HandleFunc("GET /home/", makeTemplateHandler(handlers.GetHomeHandler))
    mux.HandleFunc("GET /billing/", makeTemplateHandler(handlers.GetBillingHandler))
    mux.HandleFunc("GET /customers/", makeTemplateHandler(handlers.GetCustomersHandler))
    mux.HandleFunc("GET /products/", makeTemplateHandler(handlers.GetProductsHandler))
    mux.HandleFunc("GET /settings/", makeTemplateHandler(handlers.GetSettingsHandler))

    server := &http.Server{
        Addr: ":8080",
        Handler: mux,
    }

    log.Fatal(server.ListenAndServe())
}

func makeTemplateHandler(handlerFunction func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    //TODO: Handle refreshes (not made from htmx)
    return func(writer http.ResponseWriter, request *http.Request) {
        handlerFunction(writer, request)
    }
}

