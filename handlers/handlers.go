package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func renderTemplate(writer http.ResponseWriter, r *http.Request, templateName string, prereqTemplates []string) {
    var err error
    if r.Header.Get("HX-Request") == "true" {
        fmt.Printf("Rendering template for htmx request (content only): %v\n", templateName)
        template := getOrParse(templateName)
        err = template.ExecuteTemplate(writer, "content", nil)
    } else {
        toRender := append([]string{templateName}, prereqTemplates...)
        fmt.Printf("Rendering template for non-htmx request (entire template): %v\n", toRender)
        template := getOrParse(toRender...)
        err = template.Execute(writer, nil)
    }

    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}

var templates = make(map[string]template.Template)

func getOrParse(templateNames ...string) template.Template {
    templateKey := strings.Join(templateNames, "-")
    existing, exists := templates[templateKey]
    if exists {
        fmt.Printf("Found previously parsed template for %v\n", templateKey)
        return existing
    }

    fmt.Printf("Creating new template for %v\n", templateKey)
    templateFiles := make([]string, len(templateNames))
    for i, tmp := range templateNames {
        templateFiles[i] = "tmpl/" + tmp + ".html"
    }

    templates[templateKey] = *template.Must(template.ParseFiles(templateFiles...))
    return templates[templateKey]
}
