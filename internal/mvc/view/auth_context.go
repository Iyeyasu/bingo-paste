package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
)

// LoginContext represents a rendering context for the Login page.
type LoginContext struct {
	PageContext
}

// RegisterContext represents a rendering context for the Register page.
type RegisterContext struct {
	PageContext
}

// NewLoginContext creates a new AuthContext.
func (v *AuthView) NewLoginContext() LoginContext {
	return LoginContext{
		PageContext: PageContext{
			Page:   v.Login,
			Config: config.Get(),
		},
	}
}

// NewRegisterContext creates a new AuthContext.
func (v *AuthView) NewRegisterContext() RegisterContext {
	return RegisterContext{
		PageContext: PageContext{
			Page:   v.Register,
			Config: config.Get(),
		},
	}
}
