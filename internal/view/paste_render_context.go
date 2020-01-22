package view

import (
	"html/template"

	"github.com/Iyeyasu/bingo-paste/internal/model"
)

type pasteRenderContext struct {
	Title            string
	RawContent       string
	FormattedContent template.HTML
}

func newPasteRenderContext(paste *model.Paste) *pasteRenderContext {
	ctx := new(pasteRenderContext)
	ctx.Title = paste.Title
	ctx.RawContent = paste.RawContent
	ctx.FormattedContent = template.HTML(paste.FormattedContent)
	return ctx
}
