package view

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
)

// UserView represents the view used to render users.
type UserView struct {
	Profile *Page
	Edit    *Page
	List    *Page
}

// UserEditContext represents a rendering context for the User Edit page.
type UserEditContext struct {
	PageContext
	User *model.User
}

// UserListContext represents a rendering context for the User List page.
type UserListContext struct {
	PageContext
	Users []*model.User
}

// NewUserView creates a new UserView.
func NewUserView() *UserView {
	editPaths := []string{
		"web/template/*.go.html",
		"web/template/user/editor/*.go.html",
		"web/css/common/*.css",
	}

	listPaths := []string{
		"web/template/*.go.html",
		"web/template/user/list/*.go.html",
		"web/css/common/*.css",
	}

	v := new(UserView)
	v.Profile = NewPage("Profile", editPaths)
	v.Edit = NewPage("Edit User", editPaths)
	v.List = NewPage("List Users", listPaths)
	return v
}

// NewUserProfileContext creates a new UserEditContext.
func (v *UserView) NewUserProfileContext(r *http.Request, user *model.User) UserEditContext {
	return UserEditContext{
		User:        user,
		PageContext: NewPageContext(r, v.Profile),
	}
}

// NewUserEditContext creates a new UserEditContext.
func (v *UserView) NewUserEditContext(r *http.Request, user *model.User) UserEditContext {
	return UserEditContext{
		User:        user,
		PageContext: NewPageContext(r, v.Edit),
	}
}

// NewUserListContext creates a new UserListContext.
func (v *UserView) NewUserListContext(r *http.Request, users []*model.User) UserListContext {
	return UserListContext{
		Users:       users,
		PageContext: NewPageContext(r, v.List),
	}
}
