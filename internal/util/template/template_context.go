package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
)

// TemplateContext is the contex object passed to the html/template renderer.
type TemplateContext struct {
	View   string
	Filter string
	Config *config.Config
}

// NewTemplateContext creates a new TemplateContext.
func NewTemplateContext(view string) *TemplateContext {
	ctx := new(TemplateContext)
	ctx.View = view
	ctx.Config = config.Get()
	return ctx
}
