package view

import (
	"html/template"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
)

// ErrorView serves the view for showing errors.
type ErrorView struct {
	name     string
	template *template.Template
}

// ErrorContext contains the context for rendering the ErrorView.
type ErrorContext struct {
	template_util.TemplateContext
	StatusCode  int
	Description string
}

// NewErrorView creates a new ErrorView.
func NewErrorView() *ErrorView {
	view := new(ErrorView)
	view.name = "error"
	view.template = template_util.GetTemplate(view.name)
	return view
}

// ServeError serves the error page.
func (view *ErrorView) ServeError(w http.ResponseWriter, r *http.Request) {
	ctx := ErrorContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
		TemplateContext: template_util.TemplateContext{
			View:   view.name,
			Config: config.Get(),
		},
	}
	http_util.WriteTemplate(w, view.template, ctx)
}
