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

	view.template = template_util.GetTemplate(
		"index",
		"web/template/*.go.html",
		"web/template/error/*.go.html",
		"web/css/*.css",
	)

	return view
}

// ServeError serves the error page.
func (view *ErrorView) ServeError(w http.ResponseWriter, r *http.Request) {
	ctx := ErrorContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
		TemplateContext: template_util.TemplateContext{
			View:   "Error",
			Config: config.Get(),
		},
	}
	http_util.WriteTemplate(w, view.template, ctx)
}
