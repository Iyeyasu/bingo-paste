package view

import "net/http"

// AuthView represents the view used to render errors.
type AuthView struct {
	Login    *Page
	Register *Page
}

// LoginContext represents a rendering context for the Login page.
type LoginContext struct {
	PageContext
}

// RegisterContext represents a rendering context for the Register page.
type RegisterContext struct {
	PageContext
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
	v.Login = NewPage("Login", loginPaths)
	v.Register = NewPage("Register", registerPaths)
	return v
}

// NewLoginContext creates a new AuthContext.
func (v *AuthView) NewLoginContext(r *http.Request) LoginContext {
	return LoginContext{
		PageContext: NewPageContext(r, v.Login),
	}
}

// NewRegisterContext creates a new AuthContext.
func (v *AuthView) NewRegisterContext(r *http.Request) RegisterContext {
	return RegisterContext{
		PageContext: NewPageContext(r, v.Register),
	}
}
