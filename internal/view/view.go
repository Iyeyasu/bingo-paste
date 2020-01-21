package view

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

type view struct {
	uri      string
	template *template.Template
}

func (view *view) initialize(uri string, ctx *renderContext) {
	log.Printf("Initialize view '%s'.", uri)

	view.uri = uri
	view.template = view.renderTemplate(uri, ctx)
}

func (view *view) render(w http.ResponseWriter, req *http.Request) {
	log.Printf("Render view '%s'.", view.uri)

	w.Header().Set("Cache-Control", "public, max-age=3600")

	if req.Method != "GET" {
		log.Printf("Invalid HTTP request method %s for view %s.", req.Method, view.uri)
	}

	var paste model.Paste
	err := view.template.Execute(w, paste)
	if err != nil {
		log.Printf("Failed to render view '%s'.", view.uri)
	}
}

func (view *view) renderTemplate(name string, ctx *renderContext) *template.Template {
	log.Printf("Render view template '%s'.", name)

	var buf bytes.Buffer
	tmpl := template.Must(template.ParseGlob("web/template/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/css/*.css"))
	tmpl.ExecuteTemplate(&buf, "index.html", ctx)
	tmplMini := minifyTemplate(buf.String())

	return template.Must(template.New(name).Parse(tmplMini))
}

func minifyTemplate(str string) string {

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	res, err := m.String("text/html", str)

	if err != nil {
		log.Panicln("Failed to minify HTML template")
	}

	return res
}
