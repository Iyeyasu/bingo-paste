package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
)


type TemplateContext struct {
	Config *config.Config
}

func NewTemplateContext() *TemplateContext {
	ctx := new(TemplateContext)
	ctx.Config = config.Get()
	return ctx
}
