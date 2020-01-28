package view

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
	"github.com/julienschmidt/httprouter"
)

// EditorView .
type EditorView struct {
	name string
	html []byte
}

// NewEditorView .
func NewEditorView(store *model.PasteStore) *EditorView {
	view := new(EditorView)
	view.name = "editor"

	ctx := template_util.TemplateContext{
		View:   view.name,
		Config: config.Get(),
	}
	view.html = []byte(template_util.RenderTemplate(view.name, ctx))
	return view
}

// ServeEditor .
func (view *EditorView) ServeEditor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http_util.WriteHTML(w, view.html)
}
