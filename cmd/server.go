package main

import (
	"net/http"

	"bingo/internal/config"
	"bingo/internal/http/middleware"
	"bingo/internal/mvc/controller"
	"bingo/internal/mvc/model"
	"bingo/internal/mvc/model/store"
	"bingo/internal/session"
	"bingo/internal/util/log"

	"github.com/julienschmidt/httprouter"
)

func main() {
	db := model.NewDatabase()
	pasteStore := store.NewPasteStore(db)
	userStore := store.NewUserStore(db)
	session.Init(userStore)
	router := httprouter.New()

	errCtrl := controller.NewErrorController()
	imageCtrl := controller.NewImageController()
	pasteCtrl := controller.NewPasteController(errCtrl, pasteStore)
	userCtrl := controller.NewUserController(errCtrl, userStore)
	authCtrl := controller.NewAuthController(errCtrl, userCtrl)

	imageRoute(router, imageCtrl)
	pasteRoute(router, pasteCtrl)
	authRoute(router, authCtrl)
	userRoute(router, userCtrl)

	router.NotFound = guestMiddleware(errCtrl.ServeNotFoundError)
	log.Fatal(http.ListenAndServe(":80", router))
}

func imageRoute(router *httprouter.Router, imgCtrl *controller.ImageController) {
	router.Handler(http.MethodGet, "/favicon.ico", guestMiddleware(imgCtrl.ServeFavicon))
}

func pasteRoute(router *httprouter.Router, pasteCtrl *controller.PasteController) {
	router.Handler(http.MethodGet, "/", viewerMiddleware(pasteCtrl.ServeWritePage))
	router.Handler(http.MethodGet, "/pastes", viewerMiddleware(pasteCtrl.ServeListPage))
	router.Handler(http.MethodGet, "/pastes/:id", guestMiddleware(pasteCtrl.ServeViewPage))
	router.Handler(http.MethodGet, "/pastes/:id/raw", viewerMiddleware(pasteCtrl.ServeRawPaste))
	router.Handler(http.MethodPost, "/pastes", editorMiddleware(pasteCtrl.CreatePaste))
}

func authRoute(router *httprouter.Router, authCtrl *controller.AuthController) {
	if !config.Get().Authentication.Enabled {
		return
	}

	router.Handler(http.MethodGet, "/login", guestMiddleware(authCtrl.ServeLoginPage))
	router.Handler(http.MethodGet, "/register", guestMiddleware(authCtrl.ServeRegisterPage))
	router.Handler(http.MethodPost, "/login", guestMiddleware(authCtrl.Login))
	router.Handler(http.MethodPost, "/logout", guestMiddleware(authCtrl.Logout))

	if config.Get().Authentication.Standard.AllowRegistration {
		router.Handler(http.MethodPost, "/register", guestMiddleware(authCtrl.Register))
	}
}

func userRoute(router *httprouter.Router, userCtrl *controller.UserController) {
	if !config.Get().Authentication.Enabled {
		return
	}

	router.Handler(http.MethodGet, "/profile", viewerMiddleware(userCtrl.ServeProfilePage))
	router.Handler(http.MethodGet, "/users", adminMiddleware(userCtrl.ServeListPage))
	router.Handler(http.MethodGet, "/users/create", adminMiddleware(userCtrl.ServeCreatePage))
	router.Handler(http.MethodGet, "/users/edit/:id", adminMiddleware(userCtrl.ServeEditPage))
	router.Handler(http.MethodPost, "/profile/update", viewerMiddleware(userCtrl.UpdateProfile))
	router.Handler(http.MethodPost, "/users/create", adminMiddleware(userCtrl.CreateUser))
	router.Handler(http.MethodPost, "/users/update/:id", adminMiddleware(userCtrl.UpdateUser))
	router.Handler(http.MethodPost, "/users/delete/:id", adminMiddleware(userCtrl.DeleteUser))
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
	if !config.Get().Authentication.Enabled {
		return guestMiddleware(handler)
	}

	mw := middleware.Authorize(handler, role)
	mw = middleware.Authenticate(mw)
	return guestMiddleware(mw.ServeHTTP)
}

func guestMiddleware(handler http.HandlerFunc) http.Handler {
	mw := middleware.StartSession(handler)
	mw = middleware.TrimStrings(mw)
	mw = middleware.Log(mw)
	return mw
}
