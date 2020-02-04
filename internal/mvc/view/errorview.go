package view

import (
	"net/http"
)

// ErrorView represents the view used to render errors.
type ErrorView struct {
	Error *Page
}

// ErrorContext represents a rendering context for the Error page.
type ErrorContext struct {
	PageContext
	StatusCode int
	StatusText string
}

// NewErrorView creates a new ErrorView.
func NewErrorView() *ErrorView {
	paths := []string{
		"web/template/*.go.html",
		"web/template/error/*.go.html",
		"web/css/common/*.css",
		"web/css/error/*.css",
	}

	view := new(ErrorView)
	view.Error = NewPage("Error", "/error", paths)
	return view
}

// NewErrorContext creates a new ErrorContext.
func (v *ErrorView) NewErrorContext(r *http.Request, code int, text string) ErrorContext {
	return ErrorContext{
		StatusCode:  code,
		StatusText:  text,
		PageContext: NewPageContext(r, v.Error),
	}
}
