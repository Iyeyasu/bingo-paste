package view

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
	"github.com/julienschmidt/httprouter"
)

// ViewerView .
type ViewerView struct {
	name     string
	store    *model.PasteStore
	template *template.Template
}

// ViewerTemplateContext .
type ViewerTemplateContext struct {
	template_util.TemplateContext
	Paste *model.Paste
}

// NewViewerView .
func NewViewerView(store *model.PasteStore) *ViewerView {
	view := new(ViewerView)
	view.name = "viewer"
	view.store = store
	view.template = template_util.GetTemplate(view.name)
	return view
}

// ServePaste .
func (view *ViewerView) ServePaste(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if paste, err := view.getPaste(params); err == nil {
		ctx := ViewerTemplateContext{
			Paste: paste,
			TemplateContext: template_util.TemplateContext{
				View:   view.name,
				Config: config.Get(),
			},
		}
		http_util.WriteTemplate(w, view.template, ctx)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

// ServeRawPaste .
func (view *ViewerView) ServeRawPaste(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if paste, err := view.getPaste(params); err == nil {
		http_util.WriteText(w, []byte(paste.RawContent))
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

func (view *ViewerView) getPaste(params httprouter.Params) (*model.Paste, error) {
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return view.store.Select(id)
}
