package view

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
)

// PasteView represents the view used to render pastes.
type PasteView struct {
	Edit *Page
	List *Page
	View *Page
}

// PasteEditorContext represents a rendering context for the Paste Editor page.
type PasteEditorContext struct {
	PageContext
}

// PasteViewerContext represents a rendering context for the Paste Viewer page.
type PasteViewerContext struct {
	PageContext
	Paste *model.Paste
}

// PasteListContext represents a rendering context for the Paste List page.
type PasteListContext struct {
	PageContext
	Pastes []*model.Paste
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

// NewPasteEditorContext creates a new PasteEditorContext.
func (v *PasteView) NewPasteEditorContext(r *http.Request) PasteEditorContext {
	return PasteEditorContext{
		PageContext: NewPageContext(r, v.Edit),
	}
}

// NewPasteViewerContext creates a new PasteViewerContext.
func (v *PasteView) NewPasteViewerContext(r *http.Request, paste *model.Paste) PasteViewerContext {
	return PasteViewerContext{
		Paste:       paste,
		PageContext: NewPageContext(r, v.View),
	}
}

// NewPasteListContext creates a new PasteListContext.
func (v *PasteView) NewPasteListContext(r *http.Request, pastes []*model.Paste) PasteListContext {
	return PasteListContext{
		Pastes:      pastes,
		PageContext: NewPageContext(r, v.List),
	}
}
