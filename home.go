package main

import (
    "html/template"
    "log"
    "net/http"
    "regexp"
)
var validPath = regexp.MustCompile("^(/[a-zA-Z0-9]+/{0,1}?)*$")

// TODO: Utilize Must and only parse templates once
// var templates = template.Must(template.ParseFiles("tmpl/base.html", "tmpl/view.html", "tmpl/edit.html"))

func main() {
    http.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    http.HandleFunc("GET /", getRootHandler)
    http.HandleFunc("GET /home/", makeHandler(getHomeHandler))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getHomeHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "home")
}

func getRootHandler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/home/", 308)
}

func renderTemplate(writer http.ResponseWriter, templateName string) {
    template, err := template.ParseFiles("tmpl/base.html", "tmpl/" + templateName + ".html")
    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
    template.Execute(writer, nil)
}

func makeHandler(handlerFunction func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        path := validPath.FindStringSubmatch(request.URL.Path)
        if path == nil {
            http.NotFound(writer, request)
            return
        }

        handlerFunction(writer, request)
    }
}

