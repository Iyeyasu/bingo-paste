package main

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/http/httpmw"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/controller"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db := model.NewDatabase()
	router := new(httprouter.Router)

	// Paste routing
	pasteStore := model.NewPasteStore(db)
	pasteCtrl := controller.NewPasteController(pasteStore)
	router.Handler(http.MethodGet, "/", newMiddleware(pasteCtrl.ServeEditPage))
	router.Handler(http.MethodGet, "/pastes", newMiddleware(pasteCtrl.ServeListPage))
	router.Handler(http.MethodGet, "/pastes/:id", newMiddleware(pasteCtrl.ServeViewPage))
	router.Handler(http.MethodGet, "/pastes/:id/raw", newMiddleware(pasteCtrl.ServeRawPaste))
	router.Handler(http.MethodPost, "/pastes", newMiddleware(pasteCtrl.CreatePaste))

	// User routing
	userStore := model.NewUserStore(db)
	userCtrl := controller.NewUserController(userStore)
	router.Handler(http.MethodGet, "/users", newMiddleware(userCtrl.ServeListPage))
	router.Handler(http.MethodGet, "/users/create", newMiddleware(userCtrl.ServeCreatePage))
	router.Handler(http.MethodGet, "/users/edit/:id", newMiddleware(userCtrl.ServeEditPage))
	router.Handler(http.MethodPost, "/users/create", newMiddleware(userCtrl.CreateUser))
	router.Handler(http.MethodPost, "/users/update/:id", newMiddleware(userCtrl.UpdateUser))
	router.Handler(http.MethodPost, "/users/delete/:id", newMiddleware(userCtrl.DeleteUser))

	// Auth routing
	authCtrl := controller.NewAuthController(userStore)
	router.Handler(http.MethodGet, "/login", newMiddleware(authCtrl.ServeLoginPage))
	router.Handler(http.MethodGet, "/register", newMiddleware(authCtrl.ServeRegisterPage))
	router.Handler(http.MethodPost, "/login", newMiddleware(authCtrl.Login))
	router.Handler(http.MethodPost, "/logout", newMiddleware(authCtrl.Logout))
	router.Handler(http.MethodPost, "/register", newMiddleware(authCtrl.Register))

	// Misc controllers
	imageController := controller.NewImageController()
	errorController := controller.NewErrorController()
	router.HandlerFunc(http.MethodGet, "/favicon.ico", imageController.ServeFavicon)
	router.NotFound = newMiddleware(errorController.ServeErrorPage)

	log.Fatal(http.ListenAndServe(":80", router))
}

func newMiddleware(handler http.HandlerFunc) http.Handler {
	mw := httpmw.AuthMiddleware(handler)
	mw = httpmw.LogMiddleware(mw)
	mw = httpmw.MinifyMiddleware(mw)
	return mw
}
