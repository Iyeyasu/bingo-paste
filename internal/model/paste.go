package model

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID                 int64
	Title              string
	RawContent         string
	FormattedContent   string
	IsPublic           bool
	TimeCreatedSeconds int64
	LifetimeSeconds    int64
	Syntax             string
}

// NewPaste creates new paste.
func NewPaste() *Paste {
	paste := new(Paste)
	paste.IsPublic = false
	paste.Syntax = "auto"
	return paste
}
