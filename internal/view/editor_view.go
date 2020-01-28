package view

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
)

// EditorView serves the view for creating and viewing pastes.
type EditorView struct {
	name       string
	editorHTML []byte
}

// EditorContext contains the context for rendering the editor.
type EditorContext struct {
	template_util.TemplateContext
	URI string
}

// NewEditorView creates a new EditorView.
func NewEditorView() *EditorView {
	view := new(EditorView)
	view.name = "editor"
	view.editorHTML = view.getEditorHTML()
	return view
}

// ServeEditor serves the editor for creating pastes.
func (view *EditorView) ServeEditor(w http.ResponseWriter, r *http.Request) {
	http_util.WriteHTML(w, view.editorHTML)
}

func (view *EditorView) getEditorHTML() []byte {
	ctx := EditorContext{
		URI: "/pastes",
		TemplateContext: template_util.TemplateContext{
			View:   view.name,
			Config: config.Get(),
		},
	}
	return []byte(template_util.RenderTemplate(view.name, ctx))
}
