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
	db := model.NewDatabase()
	router := new(httprouter.Router)

	// Paste view and end point
	pasteStore := model.NewPasteStore(db)
	pasteEndPoint := api.NewPasteEndPoint(pasteStore)
	pasteView := view.NewPasteView(pasteEndPoint)
	router.Handler(http.MethodGet, "/", newMiddleware(pasteView.ServePasteEditor))
	router.Handler(http.MethodGet, "/pastes", newMiddleware(pasteView.ServePasteList))
	router.Handler(http.MethodGet, "/pastes/:id", newMiddleware(pasteView.ServePaste))
	router.Handler(http.MethodGet, "/pastes/:id/raw", newMiddleware(pasteView.ServeRawPaste))
	router.Handler(http.MethodPost, "/pastes", newMiddleware(pasteView.CreatePaste))

	// User view and end point
	userStore := model.NewUserStore(db)
	userEndPoint := api.NewUserEndPoint(userStore)
	userView := view.NewUserView(userEndPoint)
	router.Handler(http.MethodGet, "/users", newMiddleware(userView.ServeUserList))
	router.Handler(http.MethodGet, "/users/:id", newMiddleware(userView.ServeUserEditor))
	router.Handler(http.MethodPost, "/users/create", newMiddleware(userView.CreateUser))
	router.Handler(http.MethodPost, "/users/update/:id", newMiddleware(userView.UpdateUser))
	router.Handler(http.MethodPost, "/users/delete/:id", newMiddleware(userView.DeleteUser))

	// Misc views
	imageView := view.NewImageView()
	errorView := view.NewErrorView()
	router.Handler(http.MethodGet, "/favicon.ico", newMiddleware(imageView.ServeFavicon))
	router.NotFound = newMiddleware(errorView.ServeError)

	log.Fatal(http.ListenAndServe(":80", router))
}

func newMiddleware(handler http.HandlerFunc) http.Handler {
	mw := middleware.NewAuthenticationMiddleware(handler)
	mw = middleware.NewLogMiddleware(mw)
	return mw
}
