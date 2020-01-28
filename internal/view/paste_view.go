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
	name         string
	endPoint     *api.PasteEndPoint
	template     *template.Template
	listTemplate *template.Template
}

// PasteContext contains the context for rendering a single paste.
type PasteContext struct {
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
	view.name = "paste"
	view.endPoint = endPoint
	view.template = template_util.GetTemplate("viewer")
	view.listTemplate = template_util.GetTemplate("list")
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

// ServePaste serves the view for viewing individual pastes.
func (view *PasteView) ServePaste(w http.ResponseWriter, r *http.Request) {
	if paste, err := view.endPoint.GetPaste(r); err == nil {
		ctx := view.newPasteContext(paste)
		http_util.WriteTemplate(w, view.template, ctx)
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
		ctx := view.newPasteListContext(r, pastes)
		http_util.WriteTemplate(w, view.listTemplate, ctx)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

func (view *PasteView) newPasteContext(paste *model.Paste) PasteContext {
	return PasteContext{
		Paste: paste,
		TemplateContext: template_util.TemplateContext{
			View:   view.name,
			Config: config.Get(),
		},
	}
}

func (view *PasteView) newPasteListContext(r *http.Request, pastes []*model.Paste) PasteListContext {
	return PasteListContext{
		Pastes: pastes,
		TemplateContext: template_util.TemplateContext{
			View:   view.name,
			Config: config.Get(),
			Filter: r.URL.Query().Get("search"),
		},
	}
}
