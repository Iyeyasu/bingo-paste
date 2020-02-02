package view

// UserView represents the view used to render users.
type UserView struct {
	Edit *Page
	List *Page
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
	v.Edit = NewPage("Edit User", editorPaths)
	v.List = NewPage("List Users", listPaths)
	return v
}
