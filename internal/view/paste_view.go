package view

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
)

// PasteView serves the view for creating and viewing pastes.
type PasteView struct {
	endPoint       *api.PasteEndPoint
	editorTemplate *template.Template
	viewerTemplate *template.Template
	listTemplate   *template.Template
}

// PasteEditorContext contains the context for rendering the paste editor.
type PasteEditorContext struct {
	template_util.TemplateContext
	URI string
}

// PasteViewerContext contains the context for rendering a single paste.
type PasteViewerContext struct {
	template_util.TemplateContext
	Paste *model.Paste
}

// PasteListContext contains the context for rendering a list of pastes.
type PasteListContext struct {
	template_util.TemplateContext
	Pastes []*model.Paste
}

// NewPasteView creates a new PasteView.
func NewPasteView(endPoint *api.PasteEndPoint) *PasteView {
	view := new(PasteView)
	view.endPoint = endPoint

	view.editorTemplate = template_util.GetTemplate(
		"index",
		"web/template/*.go.html",
		"web/template/paste/editor/*.go.html",
		"web/css/*.css",
		"web/css/paste/*.css",
	)

	view.viewerTemplate = template_util.GetTemplate(
		"index",
		"web/template/*.go.html",
		"web/template/paste/viewer/*.go.html",
		"web/css/*.css",
		"web/css/paste/*.css",
	)

	view.listTemplate = template_util.GetTemplate(
		"index",
		"web/template/*.go.html",
		"web/template/paste/list/*.go.html",
		"web/css/*.css",
	)

	return view
}

// CreatePaste creates a new paste.
func (view *PasteView) CreatePaste(w http.ResponseWriter, r *http.Request) {
	if paste, err := view.endPoint.CreatePaste(r); err == nil {
		http.Redirect(w, r, fmt.Sprintf("/pastes/%d", paste.ID), 303)
	} else {
		http_util.WriteError(w, err.Error())
	}
}

// ServePasteEditor serves the view for creating pastes.
func (view *PasteView) ServePasteEditor(w http.ResponseWriter, r *http.Request) {
	ctx := view.newPasteEditorContext()
	http_util.WriteTemplate(w, view.editorTemplate, ctx)
}

// ServePaste serves the view for viewing individual pastes.
func (view *PasteView) ServePaste(w http.ResponseWriter, r *http.Request) {
	if paste, err := view.endPoint.GetPaste(r); err == nil {
		ctx := view.newPasteViewerContext(paste)
		http_util.WriteTemplate(w, view.viewerTemplate, ctx)
	} else {
		http_util.WriteError(w, err.Error())
	}
}

// ServeRawPaste serves the view for viewing raw text content of individual pastes.
func (view *PasteView) ServeRawPaste(w http.ResponseWriter, r *http.Request) {
	if paste, err := view.endPoint.GetPaste(r); err == nil {
		http_util.WriteText(w, []byte(paste.RawContent))
	} else {
		http_util.WriteError(w, err.Error())
	}
}

// ServePasteList serves the view for viewing a list of pastes.
func (view *PasteView) ServePasteList(w http.ResponseWriter, r *http.Request) {
	if pastes, err := view.endPoint.GetPastes(r); err == nil {
		filter := r.URL.Query().Get("search")
		ctx := view.newPasteListContext(pastes, filter)
		http_util.WriteTemplate(w, view.listTemplate, ctx)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

func (view *PasteView) newPasteEditorContext() PasteEditorContext {
	return PasteEditorContext{
		URI: "/pastes",
		TemplateContext: template_util.TemplateContext{
			View:   "Paste Editor",
			Config: config.Get(),
		},
	}
}

func (view *PasteView) newPasteViewerContext(paste *model.Paste) PasteViewerContext {
	return PasteViewerContext{
		Paste: paste,
		TemplateContext: template_util.TemplateContext{
			View:   "Paste Viewer",
			Config: config.Get(),
		},
	}
}

func (view *PasteView) newPasteListContext(pastes []*model.Paste, filter string) PasteListContext {
	return PasteListContext{
		Pastes: pastes,
		TemplateContext: template_util.TemplateContext{
			View:   "Paste List",
			Config: config.Get(),
			Filter: filter,
		},
	}
}
