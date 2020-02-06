package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/model"
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
	_, err := ctrl.login(r)
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

func (ctrl *AuthController) login(r *http.Request) (*model.User, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := ctrl.user.store.FindByUID(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = auth.CheckPasswordHash(password, user.PasswordHash.String)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errors.New("invalid username or password")
	} else if err != nil {
		return nil, err
	}

	err = session.Login(r, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
