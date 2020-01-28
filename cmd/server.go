package main

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/middleware"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/Iyeyasu/bingo-paste/internal/view"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Create object stores
	db := model.NewDatabase()
	pasteStore := model.NewPasteStore(db)

	// Route API end points
	router := new(httprouter.Router)
	pasteEndPoint := api.NewPasteEndPoint(pasteStore)
	// router.Handler(http.MethodGet, path.Join(pasteEndPoint.URI(), ":id"), test(pasteEndPoint.GetPaste))
	// router.Handler(http.MethodPost, pasteEndPoint.URI(), test(pasteEndPoint.CreatePaste))

	// Route views
	editorView := view.NewEditorView()
	errorView := view.NewErrorView()
	imageView := view.NewImageView()
	pasteView := view.NewPasteView(pasteEndPoint)
	router.Handler(http.MethodGet, "/", newMiddleware(editorView.ServeEditor))
	router.Handler(http.MethodGet, "/favicon.ico", newMiddleware(imageView.ServeFavicon))
	router.Handler(http.MethodGet, "/pastes", newMiddleware(pasteView.ServePasteList))
	router.Handler(http.MethodGet, "/pastes/:id", newMiddleware(pasteView.ServePaste))
	router.Handler(http.MethodGet, "/pastes/:id/raw", newMiddleware(pasteView.ServeRawPaste))
	router.Handler(http.MethodPost, "/pastes", newMiddleware(pasteView.CreatePaste))
	router.NotFound = newMiddleware(errorView.ServeError)

	log.Fatal(http.ListenAndServe(":80", router))
}

func newMiddleware(handler http.HandlerFunc) http.Handler {
	mw := middleware.NewAuthenticationMiddleware(handler)
	mw = middleware.NewLogMiddleware(mw)
	return mw
}
