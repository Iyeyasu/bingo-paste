package view

import (
	"net/http"

	"bingo/internal/mvc/model"
)

// PasteView represents the view used to render pastes.
type PasteView struct {
	Write *Page
	List  *Page
	View  *Page
}

// ViewPasteContext represents a rendering context for the View Paste page.
type ViewPasteContext struct {
	PageContext
	Paste *model.Paste
}

// ListPastesContext represents a rendering context for the List Pastes page.
type ListPastesContext struct {
	PageContext
	Pastes []*model.Paste
}

// NewPasteView creates a new ErrorView.
func NewPasteView() *PasteView {
	writePaths := []string{
		"web/template/*.go.html",
		"web/template/paste/write/*.go.html",
		"web/css/common/*.css",
		"web/css/paste/*.css",
	}

	listPaths := []string{
		"web/template/*.go.html",
		"web/template/paste/list/*.go.html",
		"web/css/common/*.css",
		"web/css/paste/*.css",
	}

	viewPaths := []string{
		"web/template/*.go.html",
		"web/template/paste/view/*.go.html",
		"web/css/common/*.css",
		"web/css/paste/*.css",
	}

	v := new(PasteView)
	v.Write = NewPage("Write Paste", writePaths)
	v.List = NewPage("List Pastes", listPaths)
	v.View = NewPage("View Paste", viewPaths)
	return v
}

// NewWritePasteContext creates a new PageContext for write paste page.
func (v *PasteView) NewWritePasteContext(r *http.Request) PageContext {
	return NewPageContext(r, v.Write)
}

// NewViewPasteContext creates a new PasteViewerContext.
func (v *PasteView) NewViewPasteContext(r *http.Request, paste *model.Paste) ViewPasteContext {
	return ViewPasteContext{
		Paste:       paste,
		PageContext: NewPageContext(r, v.View),
	}
}

// NewListPastesContext creates a new PasteListContext.
func (v *PasteView) NewListPastesContext(r *http.Request, pastes []*model.Paste) ListPastesContext {
	return ListPastesContext{
		Pastes:      pastes,
		PageContext: NewPageContext(r, v.List),
	}
}
