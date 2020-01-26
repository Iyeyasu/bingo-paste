package view

import (
	"html/template"
	"net/http"

	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
)

// ErrorView handles the view for creating pastes.
type ErrorView struct {
	name     string
	template *template.Template
}

// ErrorTemplateContext contains the status code and error description to render.
type ErrorTemplateContext struct {
	StatusCode  int
	Description string
}

// NewErrorView creates a new error view.
func NewErrorView() *ErrorView {
	ctx := template_util.NewTemplateContext()
	view := new(ErrorView)
	view.name = "Error"
	view.template = template_util.PrerenderTemplate("error", ctx)
	return view
}

// Serve sets up the HTTP request handling for the given URL.
func (view *ErrorView) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http_util.WriteTemplate(w, view.template, ErrorTemplateContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
	})
}
