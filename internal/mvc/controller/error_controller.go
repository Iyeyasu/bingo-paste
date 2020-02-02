package controller

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// ErrorController handles displaying errors.
type ErrorController struct {
	view *view.ErrorView
}

// NewErrorController creates a new ErrorController.
func NewErrorController() *ErrorController {
	ctrl := new(ErrorController)
	ctrl.view = view.NewErrorView()
	return ctrl
}

// ServeErrorPage serves the error page.
func (ctrl *ErrorController) ServeErrorPage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewErrorContext()
	ctrl.view.Error.Render(w, ctx)
}
