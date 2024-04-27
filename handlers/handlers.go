package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func renderTemplate(writer http.ResponseWriter, r *http.Request, templateName string, prereqTemplates []string) {
    // template, err := template.ParseFiles("tmpl/base.html", "tmpl/nav.html", "tmpl/" + templateName + ".html")
    var toRender []string

    var err error
    toRender = getTemplateFiles(templateName)
    if r.Header.Get("HX-Request") == "true" {
        fmt.Printf("Rendering template for htmx request: %v\n", toRender)
        template := template.Must(template.ParseFiles(toRender...))
        err = template.ExecuteTemplate(writer, "content", nil)
    } else {
        toRender = append(toRender, getTemplateFiles(prereqTemplates...)...)
        fmt.Printf("Rendering template for non-htmx request: %v\n", toRender)
        template := template.Must(template.ParseFiles(toRender...))
        err = template.Execute(writer, nil)
    }

    if err != nil {
        http.Error(writer, err.Error(), http.StatusInternalServerError)
    }
}

func getTemplateFiles(templates ...string) []string {
    results := make([]string, len(templates))

    for i, tmp := range templates {
        results[i] = "tmpl/" + tmp + ".html"
    }

    return results
}

