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
	StatusCode  int
	Description string
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
	view.Error = NewPage("Error", paths)
	return view
}

// NewErrorContext creates a new ErrorContext.
func (v *ErrorView) NewErrorContext(r *http.Request) ErrorContext {
	return ErrorContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
		PageContext: NewPageContext(r, v.Error),
	}
}
