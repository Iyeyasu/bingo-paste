package view

type renderContext struct {
	Title             string // Title shown in the header bar
	IconURL           string // Icon shown in the headear bar (empty for default icon)
	ReadOnly          bool   // Is the paste read-only
	EncryptionEnabled bool
	FixedRetention    bool // Do all pastes have a fixed retention period?
	FixedVisibility   bool // Do all pastes have a fixed visibility (private/public)
	JavaScript        string
}

func newRenderContext() *renderContext {
	ctx := new(renderContext)
	ctx.Title = "BINGO"
	ctx.IconURL = ""
	ctx.ReadOnly = true
	ctx.EncryptionEnabled = false
	ctx.FixedRetention = false
	ctx.FixedVisibility = false
	ctx.JavaScript = ""
	return ctx
}
