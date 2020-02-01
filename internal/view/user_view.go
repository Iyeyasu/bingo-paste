package view

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/api"
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/model"
	http_util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	template_util "github.com/Iyeyasu/bingo-paste/internal/util/template"
	"github.com/julienschmidt/httprouter"
)

// UserView serves the view for creating and viewing users.
type UserView struct {
	endPoint       *api.UserEndPoint
	editorTemplate *template.Template
	listTemplate   *template.Template
}

// UserEditorContext contains the context for rendering the user editor.
type UserEditorContext struct {
	template_util.TemplateContext
	User *model.User
}

// UserListContext contains the context for rendering a list of users.
type UserListContext struct {
	template_util.TemplateContext
	Users []*model.User
}

// NewUserView creates a new UserView.
func NewUserView(endPoint *api.UserEndPoint) *UserView {
	view := new(UserView)
	view.endPoint = endPoint

	view.editorTemplate = template_util.GetTemplate(
		"index",
		"web/template/*.go.html",
		"web/template/user/editor/*.go.html",
		"web/css/*.css",
	)

	view.listTemplate = template_util.GetTemplate(
		"index",
		"web/template/*.go.html",
		"web/template/user/list/*.go.html",
		"web/css/*.css",
	)

	return view
}

// ServeUserEditor serves the view for editing and creating a user.
func (view *UserView) ServeUserEditor(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "create" {
		ctx := view.newUserEditorContext(nil)
		http_util.WriteTemplate(w, view.editorTemplate, ctx)
	} else if user, err := view.endPoint.ReadUser(r); err == nil {
		ctx := view.newUserEditorContext(user)
		http_util.WriteTemplate(w, view.editorTemplate, ctx)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to server user editor: %s", err.Error()))
	}
}

// ServeUserList serves the view for viewing a list of users.
func (view *UserView) ServeUserList(w http.ResponseWriter, r *http.Request) {
	if users, err := view.endPoint.ReadUsers(r); err == nil {
		ctx := view.newUserListContext(users)
		http_util.WriteTemplate(w, view.listTemplate, ctx)
	} else {
		http_util.WriteError(w, fmt.Sprintf("Failed to server user list: %s", err.Error()))
	}
}

// CreateUser creates a new user.
func (view *UserView) CreateUser(w http.ResponseWriter, r *http.Request) {
	if _, err := view.endPoint.CreateUser(r); err == nil {
		http.Redirect(w, r, "/users", 303)
	} else {
		http_util.WriteError(w, err.Error())
	}
}

// UpdateUser updates an existing user.
func (view *UserView) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if _, err := view.endPoint.UpdateUser(r); err == nil {
		http.Redirect(w, r, "/users", 303)
	} else {
		http_util.WriteError(w, err.Error())
	}
}

// DeleteUser deletes an existing user.
func (view *UserView) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := view.endPoint.DeleteUser(r); err == nil {
		http.Redirect(w, r, "/users", 303)
	} else {
		http_util.WriteError(w, err.Error())
	}
}

func (view *UserView) newUserEditorContext(user *model.User) UserEditorContext {
	return UserEditorContext{
		User: user,
		TemplateContext: template_util.TemplateContext{
			View:   "User Editor",
			Config: config.Get(),
		},
	}
}

func (view *UserView) newUserListContext(users []*model.User) UserListContext {
	return UserListContext{
		Users: users,
		TemplateContext: template_util.TemplateContext{
			View:   "User List",
			Config: config.Get(),
		},
	}
}
