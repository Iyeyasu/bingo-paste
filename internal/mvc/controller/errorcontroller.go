package controller

import (
	"net/http"

	"bingo/internal/mvc/view"
)

// ErrorController handles displaying errors.
type ErrorController struct {
	errorView *view.ErrorView
}

// NewErrorController creates a new ErrorController.
func NewErrorController() *ErrorController {
	ctrl := new(ErrorController)
	ctrl.errorView = view.NewErrorView()
	return ctrl
}

// ServeNotFoundError serves a 404 not found error.
func (ctrl *ErrorController) ServeNotFoundError(w http.ResponseWriter, r *http.Request) {
	ctrl.ServeErrorPage(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// ServeUnauthorizedError serves a 401 unauthorized error.
func (ctrl *ErrorController) ServeUnauthorizedError(w http.ResponseWriter, r *http.Request) {
	ctrl.ServeErrorPage(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

// ServeInternalServerError serves a 500 internal server error.
func (ctrl *ErrorController) ServeInternalServerError(w http.ResponseWriter, r *http.Request, text string) {
	ctrl.ServeErrorPage(w, r, http.StatusInternalServerError, text)
}

// ServeErrorPage serves an error page with a custom message.
func (ctrl *ErrorController) ServeErrorPage(w http.ResponseWriter, r *http.Request, code int, text string) {
	ctx := ctrl.errorView.NewErrorContext(r, code, text)
	ctrl.errorView.Error.Render(w, ctx)
}
