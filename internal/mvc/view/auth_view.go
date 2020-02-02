package view

// AuthView represents the view used to render errors.
type AuthView struct {
	Login    *Page
	Register *Page
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
