package view

import (
	"net/http"
)

// ImageView serves images.
type ImageView struct {
}

// NewImageView creates a new ImageView.
func NewImageView() *ImageView {
	view := new(ImageView)
	return view
}

// ServeFavicon serves the favicon for the site.
func (view *ImageView) ServeFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, "web/img/favicon.ico")
}
