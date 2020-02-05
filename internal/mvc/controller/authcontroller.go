package controller

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/url"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/view"
	"bingo/internal/session"
	"bingo/internal/util/auth"
	"bingo/internal/util/log"

	"golang.org/x/crypto/bcrypt"
)

// AuthController handles user authentication.
type AuthController struct {
	err  *ErrorController
	user *UserController
	view *view.AuthView
}

// NewAuthController creates a new AuthController.
func NewAuthController(errCtrl *ErrorController, userCtrl *UserController) *AuthController {
	ctrl := new(AuthController)
	ctrl.err = errCtrl
	ctrl.user = userCtrl
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpext.ReloadWithError(w, r, "Login failed", err.Error())
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		httpext.ReloadWithError(w, r, "Login failed", err.Error())
		return
	}

	username := values.Get("username")
	password := values.Get("password")
	user, err := ctrl.user.store.FindByName(username)
	if err != nil {
		httpext.ReloadWithError(w, r, "Login failed", "Invalid username or password")
		return
	}

	err = auth.CheckPasswordHash(password, user.PasswordHash.String)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		httpext.ReloadWithError(w, r, "Login failed", "Invalid username or password")
		return
	} else if err != nil {
		httpext.ReloadWithError(w, r, "Login failed", err.Error())
		return
	}

	err = session.Login(r, user)
	if err != nil {
		httpext.ReloadWithError(w, r, "Login failed", err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out.
func (ctrl *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	err := session.Logout(r)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to log out", err.Error())
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register creates a new user.
func (ctrl *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	log.Debug("Registering a new user")

	userTmpl, err := ctrl.user.parseUserTemplate(r)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to register", err.Error())
		return
	}

	// Users aren't allowed to create their own auth settings
	authMode := config.Get().Authentication.DefaultMode
	authRole := config.Get().Authentication.DefaultRole
	theme := config.Get().Theme.Default
	userTmpl.AuthExternalID = sql.NullString{Valid: false}
	userTmpl.AuthMode = sql.NullInt32{Int32: int32(authMode), Valid: true}
	userTmpl.Role = sql.NullInt32{Int32: int32(authRole), Valid: true}
	userTmpl.Theme = sql.NullInt32{Int32: int32(theme), Valid: true}

	user, err := ctrl.user.createUser(userTmpl)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to register", err.Error())
		return
	}

	err = session.Login(r, user)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to login", err.Error())
		return
	}

	log.Debugf("User '%s' created and logged in", user.Name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
