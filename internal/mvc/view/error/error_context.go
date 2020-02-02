package view

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// ErrorContext represents a rendering context for the Error page.
type ErrorContext struct {
	view.PageContext
	StatusCode  int
	Description string
}

// NewErrorContext creates a new ErrorContext.
func (v *ErrorView) NewErrorContext() ErrorContext {
	return ErrorContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
		PageContext: view.PageContext{
			Page:   v.Error,
			Config: config.Get(),
		},
	}
}
