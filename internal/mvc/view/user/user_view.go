package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// UserView represents the view used to render users.
type UserView struct {
	Edit *view.Page
	List *view.Page
}

// NewUserView creates a new UserView.
func NewUserView() *UserView {
	editorPaths := []string{
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
	v.Edit = view.NewPage("User Editor", editorPaths)
	v.List = view.NewPage("User List", listPaths)
	return v
}
