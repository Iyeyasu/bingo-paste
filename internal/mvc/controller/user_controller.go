package controller

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
	"github.com/Iyeyasu/bingo-paste/internal/http/httpext"
)

// UserController serves the view for creating and controllering users.
type UserController struct {
	store *model.UserStore
	view  *view.UserView
}

// NewUserController creates a new UserController.
func NewUserController(store *model.UserStore) *UserController {
	ctrl := new(UserController)
	ctrl.store = store
	ctrl.view = view.NewUserView()
	return ctrl
}

// ServeCreatePage serves the view for editing and creating a user.
func (ctrl *UserController) ServeCreatePage(w http.ResponseWriter, r *http.Request) {
	ctx := ctrl.view.NewUserEditorContext(nil)
	ctrl.view.Edit.Render(w, ctx)
}

// ServeEditPage serves the view for editing and creating a user.
func (ctrl *UserController) ServeEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read user:", err.Error()))
		return
	}

	user, err := ctrl.store.FindByID(id)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read user:", err.Error()))
		return
	}

	ctx := ctrl.view.NewUserEditorContext(user)
	ctrl.view.Edit.Render(w, ctx)
}

// ServeListPage serves the view for viewing a list of users.
func (ctrl *UserController) ServeListPage(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.store.FindRange(httpext.ParseRange(r))
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to read user list", err.Error()))
		return
	}

	ctx := ctrl.view.NewUserListContext(users)
	ctrl.view.List.Render(w, ctx)
}

// CreateUser creates a new user.
func (ctrl *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	template, err := ctrl.parseUserTemplate(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to create user", err.Error()))
		return
	}

	_, err = ctrl.store.Insert(template)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to create user", err.Error()))
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// UpdateUser updates an existing user.
func (ctrl *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to update user", err.Error()))
		return
	}

	template, err := ctrl.parseUserTemplate(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to update user", err.Error()))
		return
	}
	template.ID = sql.NullInt64{Int64: id, Valid: true}

	_, err = ctrl.store.Update(template)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to update user", err.Error()))
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// DeleteUser deletes an existing user.
func (ctrl *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := httpext.ParseID(r)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to delete user:", err.Error()))
		return
	}

	err = ctrl.store.Delete(id)
	if err != nil {
		httpext.WriteError(w, fmt.Sprintln("failed to delete user:", err.Error()))
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (ctrl *UserController) parseUserTemplate(r *http.Request) (*model.UserModel, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user: %s", err)
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse user: %s", err)
	}

	user := model.UserModel{
		Password:       ctrl.parseString(values.Get("password")),
		Name:           ctrl.parseString(values.Get("name")),
		Email:          ctrl.parseString(values.Get("email")),
		AuthExternalID: ctrl.parseString(values.Get("auth_external_id")),
		AuthType:       ctrl.parseInt32(values.Get("auth_type")),
		Role:           ctrl.parseInt32(values.Get("role")),
		Theme:          ctrl.parseInt32(values.Get("theme")),
	}

	return &user, nil
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
		String: strings.TrimSpace(str),
		Valid:  str != "",
	}
}
