package controller

import (
	"net/http"
)

// ImageController handles serving images.
type ImageController struct {
}

// NewImageController creates a new ImageController.
func NewImageController() *ImageController {
	ctrl := new(ImageController)
	return ctrl
}

// ServeFavicon serves the favicon for the site.
func (ctrl *ImageController) ServeFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, "web/img/favicon.ico")
}
