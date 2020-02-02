package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// AuthView represents the view used to render errors.
type AuthView struct {
	Login    *view.Page
	Register *view.Page
}

// NewAuthView creates a new AuthView.
func NewAuthView() *AuthView {
	loginPaths := []string{
		"web/template/*.go.html",
		"web/template/auth/login/*.go.html",
		"web/css/common/*.css",
	}

	registerPaths := []string{
		"web/template/*.go.html",
		"web/template/auth/register/*.go.html",
		"web/css/common/*.css",
	}

	v := new(AuthView)
	v.Login = view.NewPage("Login", loginPaths)
	v.Register = view.NewPage("Register", registerPaths)
	return v
}
