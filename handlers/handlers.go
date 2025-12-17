package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type renderOpts struct {
	templateName    string
	prereqTemplates []string
	data            any
	entrypoints     []string
}

func newRenderOpts(name string, data any, entrypoints ...string) renderOpts {
	if len(entrypoints) == 0 {
		entrypoints = []string{"content"}
	}

	return renderOpts{
		templateName:    name,
		prereqTemplates: make([]string, 0),
		data:            data,
		entrypoints:     entrypoints,
	}
}

func renderTemplate(writer http.ResponseWriter, r *http.Request, opts renderOpts) {
	var err error

	toRender := append([]string{opts.templateName}, opts.prereqTemplates...)
	// If it's an HTMX request, we know only part of the page is being substituted.
	// Otherwise, we're rendering the whole page
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("Rendering template for htmx request (content only): %v\n", opts.templateName)
		template := getOrParse(toRender...)

		for _, ep := range opts.entrypoints {
			log.Println("Executing", ep)
			err = template.ExecuteTemplate(writer, ep, opts.data)
		}
	} else {
		toRender := append(toRender, "base")
		log.Printf("Rendering template for non-htmx request (entire template): %v\n", toRender)
		template := getOrParse(toRender...)
		err = template.Execute(writer, opts.data)
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
