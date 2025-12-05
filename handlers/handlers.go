package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func renderTemplate(writer http.ResponseWriter, r *http.Request, templateName string, prereqTemplates []string, data any) {
	var err error
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("Rendering template for htmx request (content only): %v\n", templateName)
		template := getOrParse(templateName)
		err = template.ExecuteTemplate(writer, "content", data)
	} else {
		toRender := append([]string{templateName}, prereqTemplates...)
		log.Printf("Rendering template for non-htmx request (entire template): %v\n", toRender)
		template := getOrParse(toRender...)
		err = template.Execute(writer, data)
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
		log.Printf("Found previously parsed template for %v\n", templateKey)
		return existing
	}

	log.Printf("Creating new template for %v\n", templateKey)
	templateFiles := make([]string, len(templateNames))
	for i, tmp := range templateNames {
		templateFiles[i] = "tmpl/" + tmp + ".html"
	}

	tmpl := template.New(templateNames[0] + ".html")
	addCustomFuncs(tmpl)
	template := *template.Must(tmpl.ParseFiles(templateFiles...))
	templates[templateKey] = template
	return templates[templateKey]
}

func addCustomFuncs(templ *template.Template) {
	templ.Funcs(template.FuncMap{
		"convertLocalTimestamp": func(secSinceEpoch int64) time.Time {
			return time.Unix(secSinceEpoch, 0)
		},
	})
}
