package view

import (
	"sort"
	"strings"

	"github.com/alecthomas/chroma/lexers"
)

type templateRenderContext struct {
	Title             string // Title shown in the header bar
	IconURL           string // Icon shown in the headear bar (empty for default icon)
	ReadOnly          bool   // Is the paste read-only
	UseDarkTheme      bool   // Should we used the dark theme
	EncryptionEnabled bool   // Should we enable encryption
	FixedRetention    bool   // Do all pastes have a fixed retention period?
	FixedVisibility   bool   // Do all pastes have a fixed visibility (private/public)
	Syntaxes          []string
}

func newTemplateRenderContext() *templateRenderContext {
	ctx := new(templateRenderContext)
	ctx.Title = "BINGO"
	ctx.IconURL = ""
	ctx.ReadOnly = true
	ctx.UseDarkTheme = false
	ctx.EncryptionEnabled = false
	ctx.FixedRetention = false
	ctx.FixedVisibility = false

	// Retrieve all available syntaxes.
	syntaxCount := lexers.Registry.Lexers.Len()
	ctx.Syntaxes = make([]string, 0, syntaxCount)
	for i := 0; i < syntaxCount; i++ {
		name := lexers.Registry.Lexers[i].Config().Name
		if name != "plaintext" {
			ctx.Syntaxes = append(ctx.Syntaxes, name)
		}
	}

	sort.Slice(ctx.Syntaxes, func(i, j int) bool {
		return strings.ToLower(ctx.Syntaxes[i]) < strings.ToLower(ctx.Syntaxes[j])
	})

	return ctx
}
