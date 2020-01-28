package view

import (
	"html/template"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
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
	template_util.TemplateContext
	StatusCode  int
	Description string
}

// NewErrorView creates a new error view.
func NewErrorView() *ErrorView {
	view := new(ErrorView)
	view.name = "error"
	view.template = template_util.GetTemplate(view.name)
	return view
}

// Serve sets up the HTTP request handling for the given URL.
func (view *ErrorView) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := ErrorTemplateContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
		TemplateContext: template_util.TemplateContext{
			View:   view.name,
			Config: config.Get(),
		},
	}
	http_util.WriteTemplate(w, view.template, ctx)
}
