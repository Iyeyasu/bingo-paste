package view

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/data"
)

// Renderer handles rendering of HTML templates.
// It optimizes the rendering by combining all static parts that won't change
// between HTTP requests into a single template that is then used to add in
// any variable content.
type Renderer struct {
	editorTemplate *template.Template
	viewerTemplate *template.Template
}

// Initialize pre-renders as much of the templates as possible at start-up.
func (renderer *Renderer) Initialize() {
	renderer.editorTemplate = renderEditorTemplate()
	renderer.viewerTemplate = renderViewerTemplate()
}

// RenderEditor renders the paste contents into the editable HTML.
func (renderer *Renderer) RenderEditor(rw http.ResponseWriter) {
	var paste data.Paste
	err := renderer.editorTemplate.Execute(rw, paste)
	if err != nil {
		log.Println("Failed to render editor view")
	}
}

// RenderViewer renders the paste contents into the readonly HTML.
func (renderer *Renderer) RenderViewer(rw http.ResponseWriter) {
	var paste data.Paste
	err := renderer.viewerTemplate.Execute(rw, paste)
	if err != nil {
		log.Println("Failed to render editor view")
	}
}

func renderEditorTemplate() *template.Template {
	data := getRenderData()
	data.ReadOnly = false
	return render(&data)
}

func renderViewerTemplate() *template.Template {
	data := getRenderData()
	data.ReadOnly = true
	return render(&data)
}

func render(data *RenderData) *template.Template {
	// Add all non changing content to the template
	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(&buf, data)

	// Create a new template that only needs to be filled with changing content
	return template.Must(template.New("View").Parse(buf.String()))
}

func getRenderData() RenderData {
	var data RenderData
	return data
}
