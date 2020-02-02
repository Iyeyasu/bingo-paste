package controller

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	model "github.com/Iyeyasu/bingo-paste/internal/mvc/model/user"
	view "github.com/Iyeyasu/bingo-paste/internal/mvc/view/auth"
	"github.com/Iyeyasu/bingo-paste/internal/util"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

// AuthController handles user authentication.
type AuthController struct {
	store *model.UserStore
	view  *view.AuthView
}

// NewAuthController creates a new AuthController.
func NewAuthController(store *model.UserStore) *AuthController {
	ctrl := new(AuthController)
	ctrl.store = store
	ctrl.view = view.NewAuthView()
	return ctrl
}

// ServeLoginPage serves the login page.
func (ctrl *AuthController) ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewLoginContext()
	ctrl.view.Login.Render(w, ctx)
}

// ServeRegisterPage serves the register page.
func (ctrl *AuthController) ServeRegisterPage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewRegisterContext()
	ctrl.view.Register.Render(w, ctx)
}

// Login authenticates a user.
func (ctrl *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	log.Debug("Logging in")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to parse login info:", err))
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to parse login info:", err))
		return
	}

	username := strings.TrimSpace(values.Get("username"))
	password := strings.TrimSpace(values.Get("password"))
	log.Debugf("Login username '%s", username)
	log.Tracef("Login password '%s'", password)

	user, err := ctrl.store.FindByName(username)
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to login user:", err.Error()))
		return
	}

	err = util.CheckPasswordHash(password, user.PasswordHash)
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to login user:", err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out.
func (ctrl *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Register creates a new user.
func (ctrl *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	log.Debug("Registering a new user")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to parse registering info:", err))
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to parse registering info:", err))
		return
	}

	username := strings.TrimSpace(values.Get("username"))
	email := strings.TrimSpace(values.Get("email"))
	password := strings.TrimSpace(values.Get("password"))
	passwordCheck := strings.TrimSpace(values.Get("password_confirm"))

	log.Debugf("Registering username '%s", username)
	log.Debugf("Registering email '%s", email)
	log.Tracef("Registering password '%s' ('%s')", password, passwordCheck)

	if password != passwordCheck {
		http_util.WriteError(w, fmt.Sprint("failed to register user: passwords don't match"))
		return
	}

	template := model.UserModel{
		Password:       sql.NullString{String: password, Valid: true},
		Name:           sql.NullString{String: username, Valid: true},
		Email:          sql.NullString{String: email, Valid: true},
		AuthType:       sql.NullInt32{Int32: int32(model.AuthStandard), Valid: true},
		AuthExternalID: sql.NullString{Valid: false},
		Role:           sql.NullInt32{Int32: int32(model.RoleEditor), Valid: true},
		Theme:          sql.NullInt32{Int32: int32(model.ThemeLight), Valid: true},
	}

	_, err = ctrl.store.Insert(&template)
	if err != nil {
		http_util.WriteError(w, fmt.Sprintln("failed to create user", err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
