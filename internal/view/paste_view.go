package view

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
	"github.com/julienschmidt/httprouter"
)

var (
	editorTemplate = "editor.html"
	viewerTemplate = "viewer.html"
)

// PasteView .
type PasteView struct {
	store          *model.PasteStore
	editorHTML     []byte
	viewerTemplate *template.Template
}

// NewPasteView .
func NewPasteView(store *model.PasteStore) *PasteView {
	ctx := template_util.NewTemplateContext()
	view := new(PasteView)
	view.store = store
	view.editorHTML = []byte(template_util.RenderTemplate("editor", ctx))
	view.viewerTemplate = template_util.PrerenderTemplate("viewer", ctx)
	return view
}

// ServeEditor .
func (view *PasteView) ServeEditor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http_util.WriteHTML(w, view.editorHTML)
}

// ServePaste .
func (view *PasteView) ServePaste(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if paste, err := view.getPaste(params); err == nil {
		http_util.WriteTemplate(w, view.viewerTemplate, paste)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

// ServeRawPaste .
func (view *PasteView) ServeRawPaste(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if paste, err := view.getPaste(params); err == nil {
		http_util.WriteText(w, []byte(paste.RawContent))
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to retrieve paste: %s", err.Error()))
	}
}

func (view *PasteView) getPaste(params httprouter.Params) (*model.Paste, error) {
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return view.store.Select(id)
}
