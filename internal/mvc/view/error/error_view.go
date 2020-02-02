package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// ErrorView represents the view used to render errors.
type ErrorView struct {
	Error *view.Page
}

// NewErrorView creates a new ErrorView.
func NewErrorView() *ErrorView {
	paths := []string{
		"web/template/*.go.html",
		"web/template/error/*.go.html",
		"web/css/common/*.css",
		"web/css/error/*.css",
	}

	v := new(ErrorView)
	v.Error = view.NewPage("Error", paths)
	return v
}
