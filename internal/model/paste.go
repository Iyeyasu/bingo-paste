package model

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID                 int64
	Title              string
	Content            string
	IsPublic           bool
	TimeCreatedSeconds int64
	LifetimeSeconds    int64
	Syntax             string
}

// NewPaste creates new paste.
func NewPaste() *Paste {
	paste := new(Paste)
	paste.ID = 0
	paste.Title = ""
	paste.Content = ""
	paste.IsPublic = false
	paste.TimeCreatedSeconds = 0
	paste.LifetimeSeconds = 0
	paste.Syntax = "auto"
	return paste
}
