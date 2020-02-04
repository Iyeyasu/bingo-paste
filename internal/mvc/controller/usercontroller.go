package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/mvc/model"
	"bingo/internal/mvc/model/store"
	"bingo/internal/mvc/view"
	"bingo/internal/session"
	"bingo/internal/util/log"
)

// UserController serves the view for creating and controllering users.
type UserController struct {
	store *store.UserStore
	view  *view.UserView
}

// NewUserController creates a new UserController.
func NewUserController(store *store.UserStore) *UserController {
	ctrl := new(UserController)
	ctrl.store = store
	ctrl.view = view.NewUserView()

	if ctrl.store.Count() == 0 {
		ctrl.createInitialUser()
	}

	return ctrl
}

// ServeCreatePage serves the view for editing and creating a user.
func (ctrl *UserController) ServeCreatePage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewEditUserContext(r, nil)
	ctrl.view.Edit.Render(w, ctx)
}

// ServeProfilePage serves the view for editing users own profile.
func (ctrl *UserController) ServeProfilePage(w http.ResponseWriter, r *http.Request) {
	user := session.User(r)
	if user == nil {
		httpext.InternalError(w, "failed to serve profile page: no user authenticated")
		return
	}

	ctx := ctrl.view.NewEditProfileContext(r, user)
	ctrl.view.Profile.Render(w, ctx)
}

// ServeEditPage serves the view for editing and creating a user.
func (ctrl *UserController) ServeEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to serve edit user page:", err.Error()))
		return
	}

	user, err := ctrl.store.FindByID(id)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to serve edit user page:", err.Error()))
		return
	}

	ctx := ctrl.view.NewEditUserContext(r, user)
	ctrl.view.Edit.Render(w, ctx)
}

// ServeListPage serves the view for viewing a list of users.
func (ctrl *UserController) ServeListPage(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.store.FindRange(httpext.ParseRange(r))
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to serve list users page", err.Error()))
		return
	}

	ctx := ctrl.view.NewListUsersContext(r, users)
	ctrl.view.List.Render(w, ctx)
}

// CreateUser creates a new user.
func (ctrl *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	template, err := ctrl.parseUserTemplate(r)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to create user", err.Error()))
		return
	}

	_, err = ctrl.store.Insert(template)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to create user", err.Error()))
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// UpdateProfile updates the profile of the logged in user.
func (ctrl *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userTmpl, err := ctrl.parseProfileTemplate(r)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to update profile", err.Error()))
		return
	}

	_, err = ctrl.store.Update(userTmpl)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to update profile", err.Error()))
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// UpdateUser updates an existing user.
func (ctrl *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userTmpl, err := ctrl.parseUserTemplate(r)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to update user", err.Error()))
		return
	}

	_, err = ctrl.store.Update(userTmpl)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to update user", err.Error()))
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// DeleteUser deletes an existing user.
func (ctrl *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to delete user:", err.Error()))
		return
	}

	user := session.User(r)
	if user != nil && id == user.ID {
		httpext.InternalError(w, "failed to delete user: can't delete self")
		return
	}

	err = ctrl.store.Delete(id)
	if err != nil {
		httpext.InternalError(w, fmt.Sprintln("failed to delete user:", err.Error()))
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (ctrl *UserController) createInitialUser() error {
	log.Debug("Creating initial admin user")

	userTmpl := model.UserTemplate{
		Password:       sql.NullString{String: "admin", Valid: true},
		Name:           sql.NullString{String: "admin", Valid: true},
		Email:          sql.NullString{String: "admin@localhost", Valid: true},
		AuthExternalID: sql.NullString{String: "", Valid: false},
		AuthMode:       sql.NullInt32{Int32: int32(config.AuthStandard), Valid: true},
		Role:           sql.NullInt32{Int32: int32(config.RoleAdmin), Valid: true},
		Theme:          sql.NullInt32{Int32: int32(config.ThemeLight), Valid: true},
	}

	_, err := ctrl.store.Insert(&userTmpl)
	return err
}

func (ctrl *UserController) parseProfileTemplate(r *http.Request) (*model.UserTemplate, error) {
	user := session.User(r)
	if user == nil {
		return nil, errors.New("user is nil")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	userTmpl := model.UserTemplate{
		ID:       sql.NullInt64{Int64: user.ID, Valid: true},
		Password: ctrl.parseString(values.Get("password")),
		Name:     ctrl.parseString(values.Get("name")),
		Email:    ctrl.parseString(values.Get("email")),
		Theme:    ctrl.parseInt32(values.Get("theme")),
	}

	return &userTmpl, nil
}

func (ctrl *UserController) parseUserTemplate(r *http.Request) (*model.UserTemplate, error) {
	id, err := httpext.ParseID(r)
	userID := sql.NullInt64{Int64: id, Valid: true}
	if err != nil {
		userID.Int64 = 0
		userID.Valid = false
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	userTmpl := model.UserTemplate{
		ID:             userID,
		Password:       ctrl.parseString(values.Get("password")),
		Name:           ctrl.parseString(values.Get("name")),
		Email:          ctrl.parseString(values.Get("email")),
		AuthExternalID: ctrl.parseString(values.Get("auth_external_id")),
		AuthMode:       ctrl.parseInt32(values.Get("auth_mode")),
		Role:           ctrl.parseInt32(values.Get("role")),
		Theme:          ctrl.parseInt32(values.Get("theme")),
	}

	return &userTmpl, nil
}

func (ctrl *UserController) parseInt32(str string) sql.NullInt32 {
	val, err := strconv.Atoi(str)
	return sql.NullInt32{
		Int32: int32(val),
		Valid: err == nil,
	}
}

func (ctrl *UserController) parseString(str string) sql.NullString {
	return sql.NullString{
		String: str,
		Valid:  str != "",
	}
}
