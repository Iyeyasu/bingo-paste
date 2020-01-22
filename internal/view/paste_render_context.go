package view

import (
	"html/template"

	"github.com/Iyeyasu/bingo-paste/internal/model"
)

type pasteRenderContext struct {
	Title   string
	Content template.HTML
}

func newPasteRenderContext(paste *model.Paste) *pasteRenderContext {
	ctx := new(pasteRenderContext)
	ctx.Title = paste.Title
	ctx.Content = template.HTML(paste.Content)
	return ctx
}
