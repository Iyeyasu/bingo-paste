package main

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/http/middleware"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/controller"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model/store"
	"github.com/Iyeyasu/bingo-paste/internal/session"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
	"github.com/julienschmidt/httprouter"
)

func main() {
	db := model.NewDatabase()
	pasteStore := store.NewPasteStore(db)
	userStore := store.NewUserStore(db)
	router := httprouter.New()
	session.InitDefault(userStore)

	// Paste routing
	pasteCtrl := controller.NewPasteController(pasteStore)
	router.Handler(http.MethodGet, "/", viewerMiddleware(pasteCtrl.ServeEditPage))
	router.Handler(http.MethodGet, "/pastes", viewerMiddleware(pasteCtrl.ServeListPage))
	router.Handler(http.MethodGet, "/pastes/:id", guestMiddleware(pasteCtrl.ServeViewPage))
	router.Handler(http.MethodGet, "/pastes/:id/raw", viewerMiddleware(pasteCtrl.ServeRawPaste))
	router.Handler(http.MethodPost, "/pastes", editorMiddleware(pasteCtrl.CreatePaste))

	// User routing
	userCtrl := controller.NewUserController(userStore)
	router.Handler(http.MethodGet, "/profile", viewerMiddleware(userCtrl.ServeProfilePage))
	router.Handler(http.MethodGet, "/users", adminMiddleware(userCtrl.ServeListPage))
	router.Handler(http.MethodGet, "/users/create", adminMiddleware(userCtrl.ServeCreatePage))
	router.Handler(http.MethodGet, "/users/edit/:id", adminMiddleware(userCtrl.ServeEditPage))
	router.Handler(http.MethodPost, "/profile/update", viewerMiddleware(userCtrl.UpdateProfile))
	router.Handler(http.MethodPost, "/users/create", adminMiddleware(userCtrl.CreateUser))
	router.Handler(http.MethodPost, "/users/update/:id", adminMiddleware(userCtrl.UpdateUser))
	router.Handler(http.MethodPost, "/users/delete/:id", adminMiddleware(userCtrl.DeleteUser))

	// Auth routing
	authCtrl := controller.NewAuthController(userStore)
	router.Handler(http.MethodGet, "/login", guestMiddleware(authCtrl.ServeLoginPage))
	router.Handler(http.MethodGet, "/register", guestMiddleware(authCtrl.ServeRegisterPage))
	router.Handler(http.MethodPost, "/login", guestMiddleware(authCtrl.Login))
	router.Handler(http.MethodPost, "/logout", guestMiddleware(authCtrl.Logout))
	router.Handler(http.MethodPost, "/register", guestMiddleware(authCtrl.Register))

	// Misc routing
	imageController := controller.NewImageController()
	errorController := controller.NewErrorController()
	router.Handler(http.MethodGet, "/favicon.ico", guestMiddleware(imageController.ServeFavicon))
	router.NotFound = guestMiddleware(errorController.ServeErrorPage)

	log.Fatal(http.ListenAndServe(":80", router))
}

func adminMiddleware(handler http.HandlerFunc) http.Handler {
	return authMiddleware(handler, config.RoleAdmin)
}

func editorMiddleware(handler http.HandlerFunc) http.Handler {
	return authMiddleware(handler, config.RoleEditor)
}

func viewerMiddleware(handler http.HandlerFunc) http.Handler {
	return authMiddleware(handler, config.RoleViewer)
}

func authMiddleware(handler http.HandlerFunc, role config.Role) http.Handler {
	mw := middleware.Authorize(handler, role)
	mw = middleware.Authenticate(mw)
	mw = guestMiddleware(mw.ServeHTTP)
	return mw
}

func guestMiddleware(handler http.HandlerFunc) http.Handler {
	mw := middleware.StartSession(handler)
	mw = middleware.TrimStrings(mw)
	mw = middleware.Minify(mw)
	mw = middleware.Log(mw)
	return mw
}
