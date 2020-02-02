package view

// PasteView represents the view used to render pastes.
type PasteView struct {
	Edit *Page
	List *Page
	View *Page
}

// NewPasteView creates a new ErrorView.
func NewPasteView() *PasteView {
	editorPaths := []string{
		"web/template/*.go.html",
		"web/template/paste/editor/*.go.html",
		"web/css/common/*.css",
		"web/css/paste/*.css",
	}

	listPaths := []string{
		"web/template/*.go.html",
		"web/template/paste/list/*.go.html",
		"web/css/common/*.css",
		"web/css/paste/*.css",
	}

	viewerPaths := []string{
		"web/template/*.go.html",
		"web/template/paste/viewer/*.go.html",
		"web/css/common/*.css",
		"web/css/paste/*.css",
	}

	v := new(PasteView)
	v.Edit = NewPage("Create Paste", editorPaths)
	v.List = NewPage("List Pastes", listPaths)
	v.View = NewPage("View Paste", viewerPaths)
	return v
}
