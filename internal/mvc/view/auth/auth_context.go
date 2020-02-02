package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/view"
)

// LoginContext represents a rendering context for the Login page.
type LoginContext struct {
	view.PageContext
}

// RegisterContext represents a rendering context for the Register page.
type RegisterContext struct {
	view.PageContext
}

// NewLoginContext creates a new AuthContext.
func (v *AuthView) NewLoginContext() LoginContext {
	return LoginContext{
		PageContext: view.PageContext{
			Page:   v.Login,
			Config: config.Get(),
		},
	}
}

// NewRegisterContext creates a new AuthContext.
func (v *AuthView) NewRegisterContext() RegisterContext {
	return RegisterContext{
		PageContext: view.PageContext{
			Page:   v.Register,
			Config: config.Get(),
		},
	}
}
