package view

// ErrorView represents the view used to render errors.
type ErrorView struct {
	Error *Page
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
