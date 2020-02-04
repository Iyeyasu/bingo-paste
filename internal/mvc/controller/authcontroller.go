package controller

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/model"
	"bingo/internal/mvc/model/store"
	"bingo/internal/mvc/view"
	"bingo/internal/session"
	"bingo/internal/util/auth"
	"bingo/internal/util/log"
)

// AuthController handles user authentication.
type AuthController struct {
	store *store.UserStore
	view  *view.AuthView
}

// NewAuthController creates a new AuthController.
func NewAuthController(store *store.UserStore) *AuthController {
	ctrl := new(AuthController)
	ctrl.store = store
	ctrl.view = view.NewAuthView()
	return ctrl
}

// ServeLoginPage serves the login page.
func (ctrl *AuthController) ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewLoginContext(r)
	ctrl.view.Login.Render(w, ctx)
}

// ServeRegisterPage serves the register page.
func (ctrl *AuthController) ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewRegisterContext(r)
	ctrl.view.Register.Render(w, ctx)
}

// Login authenticates a user.
func (ctrl *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("Logging in")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to parse login info:", err))
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to parse login info:", err))
		return
	}

	username := values.Get("username")
	password := values.Get("password")
	log.Debugf("Login username '%s", username)
	log.Tracef("Login password '%s'", password)

	user, err := ctrl.store.FindByName(username)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to login user:", err.Error()))
		return
	}

	err = auth.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to login user:", err.Error()))
		return
	}

	err = session.Login(r, user)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to login user:", err.Error()))
		return
	}

	log.Debugf("User '%s' logged in", username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out.
func (ctrl *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	session.Logout(r)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register creates a new user.
func (ctrl *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	log.Debug("Registering a new user")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to parse registering info:", err))
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to parse registering info:", err))
		return
	}

	username := values.Get("username")
	email := values.Get("email")
	password := values.Get("password")
	passwordCheck := values.Get("password_confirm")

	log.Debugf("Registering username '%s", username)
	log.Debugf("Registering email '%s", email)
	log.Tracef("Registering password '%s' ('%s')", password, passwordCheck)

	if password != passwordCheck {
		httpext.InternalError(w, fmt.Sprint("failed to register user: passwords don't match"))
		return
	}

	theme := config.Get().Theme.Default
	authRole := config.Get().Authentication.DefaultRole
	authMode := config.Get().Authentication.DefaultMode
	userTmpl := model.UserTemplate{
		Password:       sql.NullString{String: password, Valid: true},
		Name:           sql.NullString{String: username, Valid: true},
		Email:          sql.NullString{String: email, Valid: true},
		AuthMode:       sql.NullInt32{Int32: int32(authMode), Valid: true},
		AuthExternalID: sql.NullString{Valid: false},
		Role:           sql.NullInt32{Int32: int32(authRole), Valid: true},
		Theme:          sql.NullInt32{Int32: int32(theme), Valid: true},
	}

	user, err := ctrl.store.Insert(&userTmpl)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to create user", err.Error()))
		return
	}

	err = session.Login(r, user)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to login created user:", err.Error()))
		return
	}

	log.Debugf("User '%s' created and logged in", username)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
