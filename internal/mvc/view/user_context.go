package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
)

// UserEditorContext represents a rendering context for the User Editor page.
type UserEditorContext struct {
	PageContext
	User *model.User
}

// UserListContext represents a rendering context for the User List page.
type UserListContext struct {
	PageContext
	Users []*model.User
}

// NewUserEditorContext creates a new UserEditorContext.
func (v *UserView) NewUserEditorContext(user *model.User) UserEditorContext {
	return UserEditorContext{
		User: user,
		PageContext: PageContext{
			Page:   v.Edit,
			Config: config.Get(),
		},
	}
}

// NewUserListContext creates a new UserListContext.
func (v *UserView) NewUserListContext(users []*model.User) UserListContext {
	return UserListContext{
		Users: users,
		PageContext: PageContext{
			Page:   v.List,
			Config: config.Get(),
		},
	}
}
