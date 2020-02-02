package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// PasteView represents the view used to render pastes.
type PasteView struct {
	Edit *view.Page
	List *view.Page
	View *view.Page
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
	v.Edit = view.NewPage("Edit Paste", editorPaths)
	v.List = view.NewPage("List Pastes", listPaths)
	v.View = view.NewPage("View Paste", viewerPaths)
	return v
}
