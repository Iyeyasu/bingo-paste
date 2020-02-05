package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
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
	err   *ErrorController
	store *store.UserStore
	view  *view.UserView
}

// NewUserController creates a new UserController.
func NewUserController(errCtrl *ErrorController, store *store.UserStore) *UserController {
	ctrl := new(UserController)
	ctrl.err = errCtrl
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
		ctrl.err.ServeInternalServerError(w, r, "Failed to serve profile page: no user authenticated")
		return
	}

	ctx := ctrl.view.NewEditProfileContext(r, user)
	ctrl.view.Profile.Render(w, ctx)
}

// ServeEditPage serves the view for editing and creating a user.
func (ctrl *UserController) ServeEditPage(w http.ResponseWriter, r *http.Request) {
	user, err := ctrl.getUser(r)
	if err != nil {
		ctrl.err.ServeInternalServerError(w, r, fmt.Sprintln("Failed to serve edit user page:", err.Error()))
		return
	}

	ctx := ctrl.view.NewEditUserContext(r, user)
	ctrl.view.Edit.Render(w, ctx)
}

// ServeListPage serves the view for viewing a list of users.
func (ctrl *UserController) ServeListPage(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.store.FindRange(httpext.ParseRange(r))
	if err != nil {
		ctrl.err.ServeInternalServerError(w, r, fmt.Sprintln("Failed to serve list users page:", err.Error()))
		return
	}

	ctx := ctrl.view.NewListUsersContext(r, users)
	ctrl.view.List.Render(w, ctx)
}

// CreateUser creates a new user.
func (ctrl *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	userTmpl, err := ctrl.parseUserTemplate(r)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to create user", err.Error())
		return
	}

	user, err := ctrl.createUser(userTmpl)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to create user", err.Error())
		return
	}

	msg := fmt.Sprintf("User <b>%s</b> created successfully", user.Name)
	note := model.NewSuccessNotification("Created", msg)
	httpext.RedirectWithNotify(w, r, "/users", http.StatusSeeOther, note)
}

// UpdateProfile updates the profile of the logged in user.
func (ctrl *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	_, err := ctrl.updateProfile(r)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to save profile", err.Error())
		return
	}

	msg := "Profile saved successfully"
	note := model.NewSuccessNotification("Saved", msg)
	httpext.RedirectWithNotify(w, r, "/profile", http.StatusSeeOther, note)
}

// UpdateUser updates an existing user.
func (ctrl *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, err := ctrl.updateUser(r)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to save user", err.Error())
		return
	}

	msg := fmt.Sprintf("User <b>%s</b> saved successfully", user.Name)
	note := model.NewSuccessNotification("Saved", msg)
	httpext.RedirectWithNotify(w, r, "/users", http.StatusSeeOther, note)
}

// DeleteUser deletes an existing user.
func (ctrl *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := ctrl.deleteUser(r)
	if err != nil {
		httpext.ReloadWithError(w, r, "Failed to delete user", err.Error())
		return
	}

	msg := fmt.Sprintf("User <b>%s</b> deleted successfully", user.Name)
	note := model.NewSuccessNotification("Deleted", msg)
	httpext.RedirectWithNotify(w, r, "/users", http.StatusSeeOther, note)
}

func (ctrl *UserController) getUser(r *http.Request) (*model.User, error) {
	id, err := httpext.ParseID(r)
	if err != nil {
		return nil, err
	}

	return ctrl.store.FindByID(id)
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

func (ctrl *UserController) createUser(userTmpl *model.UserTemplate) (*model.User, error) {
	_, err := ctrl.store.FindByName(userTmpl.Name.String)
	if err != sql.ErrNoRows {
		return nil, errors.New("username already taken")
	}

	user, err := ctrl.store.Insert(userTmpl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes an existing user.
func (ctrl *UserController) deleteUser(r *http.Request) (*model.User, error) {
	id, err := httpext.ParseID(r)
	if err != nil {
		return nil, err
	}

	currentUser := session.User(r)
	if currentUser != nil && id == currentUser.ID {
		return nil, errors.New("can't delete self")
	}

	user, err := ctrl.store.FindByID(id)
	if err != nil {
		return nil, err
	}

	err = ctrl.store.Delete(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ctrl *UserController) updateUser(r *http.Request) (*model.User, error) {
	userTmpl, err := ctrl.parseUserTemplate(r)
	if err != nil {
		return nil, err
	}

	user, err := ctrl.store.Update(userTmpl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ctrl *UserController) updateProfile(r *http.Request) (*model.User, error) {
	userTmpl, err := ctrl.parseProfileTemplate(r)
	if err != nil {
		return nil, err
	}

	user, err := ctrl.store.Update(userTmpl)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ctrl *UserController) parseProfileTemplate(r *http.Request) (*model.UserTemplate, error) {
	userTmpl, err := ctrl.parseUserTemplate(r)
	if err != nil {
		return nil, err
	}

	user := session.User(r)
	if user == nil {
		return nil, errors.New("no user logged in")
	}

	// Users aren't allowed to change their own auth settings
	userTmpl.ID = sql.NullInt64{Int64: user.ID, Valid: true}
	userTmpl.AuthMode = sql.NullInt32{Valid: false}
	userTmpl.AuthExternalID = sql.NullString{Valid: false}
	userTmpl.Role = sql.NullInt32{Valid: false}

	return userTmpl, nil
}

func (ctrl *UserController) parseUserTemplate(r *http.Request) (*model.UserTemplate, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	password := r.FormValue("password")
	if passwordCheck, ok := r.Form["password_confirm"]; ok && password != passwordCheck[0] {
		return nil, errors.New("passwords do not match")
	}

	id, err := httpext.ParseID(r)
	userID := sql.NullInt64{Int64: id, Valid: true}
	if err != nil {
		userID.Int64 = 0
		userID.Valid = false
	}

	userTmpl := model.UserTemplate{
		ID:             userID,
		Password:       ctrl.parseString(r.FormValue("password")),
		Name:           ctrl.parseString(r.FormValue("username")),
		Email:          ctrl.parseString(r.FormValue("email")),
		AuthExternalID: ctrl.parseString(r.FormValue("auth_external_id")),
		AuthMode:       ctrl.parseInt32(r.FormValue("auth_mode")),
		Role:           ctrl.parseInt32(r.FormValue("role")),
		Theme:          ctrl.parseInt32(r.FormValue("theme")),
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
