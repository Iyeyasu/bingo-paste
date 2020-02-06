package view

import (
	"net/http"

	"bingo/internal/mvc/model"
)

// UserView represents the view used to render users.
type UserView struct {
	Profile *Page
	Edit    *Page
	List    *Page
}

// EditUserContext represents a rendering context for the User Edit page.
type EditUserContext struct {
	PageContext
	User *model.User
}

// ListUsersContext represents a rendering context for the User List page.
type ListUsersContext struct {
	PageContext
	TotalCount int
	Users      []model.User
}

// NewUserView creates a new UserView.
func NewUserView() *UserView {
	editPaths := []string{
		"web/template/*.go.html",
		"web/template/user/edit/*.go.html",
		"web/css/common/*.css",
	}

	listPaths := []string{
		"web/template/*.go.html",
		"web/template/user/list/*.go.html",
		"web/css/common/*.css",
	}

	v := new(UserView)
	v.Profile = NewPage("Profile", "/profile", editPaths)
	v.Edit = NewPage("Edit User", "users/:id", editPaths)
	v.List = NewPage("List Users", "users", listPaths)
	return v
}

// NewEditProfileContext creates a new EditUserContext.
func (v *UserView) NewEditProfileContext(r *http.Request, user *model.User) EditUserContext {
	return EditUserContext{
		User:        user,
		PageContext: NewPageContext(r, v.Profile),
	}
}

// NewEditUserContext creates a new EditUserContext.
func (v *UserView) NewEditUserContext(r *http.Request, user *model.User) EditUserContext {
	return EditUserContext{
		User:        user,
		PageContext: NewPageContext(r, v.Edit),
	}
}

// NewListUsersContext creates a new ListUsersContext.
func (v *UserView) NewListUsersContext(r *http.Request, users []model.User) ListUsersContext {
	return ListUsersContext{
		Users:       users,
		PageContext: NewPageContext(r, v.List),
	}
}
