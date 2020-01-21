package view

import (
	"net/http"
)

// CopyView serves the viewer page of the application.
type CopyView struct {
	view
}

// Handle sets up the HTTP request handling for the given URI.
func (view *CopyView) Handle(uri string) {
	ctx := newRenderContext()
	ctx.ReadOnly = true
	view.initialize(uri, ctx)

	http.HandleFunc(uri, view.render)
}
