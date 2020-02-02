package view

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
)

// ErrorContext represents a rendering context for the Error page.
type ErrorContext struct {
	PageContext
	StatusCode  int
	Description string
}

// NewErrorContext creates a new ErrorContext.
func (v *ErrorView) NewErrorContext() ErrorContext {
	return ErrorContext{
		StatusCode:  http.StatusNotFound,
		Description: "Page not found",
		PageContext: PageContext{
			Page:   v.Error,
			Config: config.Get(),
		},
	}
}
