package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
)

// PageContext represents a rendering context for a page template.
type PageContext struct {
	Page        *Page
	CurrentUser *model.User
	Filter      string
	Config      *config.Config
}
