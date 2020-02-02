package view

import (
	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
)

// PasteEditorContext represents a rendering context for the Paste Editor page.
type PasteEditorContext struct {
	PageContext
}

// PasteViewerContext represents a rendering context for the Paste Viewer page.
type PasteViewerContext struct {
	PageContext
	Paste *model.Paste
}

// PasteListContext represents a rendering context for the Paste List page.
type PasteListContext struct {
	PageContext
	Pastes []*model.Paste
}

// NewPasteEditorContext creates a new PasteEditorContext.
func (v *PasteView) NewPasteEditorContext() PasteEditorContext {
	return PasteEditorContext{
		PageContext: PageContext{
			Page:   v.Edit,
			Config: config.Get(),
		},
	}
}

// NewPasteViewerContext creates a new PasteViewerContext.
func (v *PasteView) NewPasteViewerContext(user *model.Paste) PasteViewerContext {
	return PasteViewerContext{
		Paste: user,
		PageContext: PageContext{
			Page:   v.View,
			Config: config.Get(),
		},
	}
}

// NewPasteListContext creates a new PasteListContext.
func (v *PasteView) NewPasteListContext(users []*model.Paste, filter string) PasteListContext {
	return PasteListContext{
		Pastes: users,
		PageContext: PageContext{
			Page:   v.List,
			Filter: filter,
			Config: config.Get(),
		},
	}
}
