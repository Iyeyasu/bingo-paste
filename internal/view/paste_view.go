package view

import (
	"net/http"
)

// PasteView serves the editor page of the application.
type PasteView struct {
	view
}

// Handle sets up the HTTP request handling for the given URI.
func (view *PasteView) Handle(uri string) {
	ctx := newRenderContext()
	ctx.ReadOnly = false
	view.initialize(uri, ctx)

	http.HandleFunc(uri, view.render)
}
