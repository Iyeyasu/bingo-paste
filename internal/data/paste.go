package data

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	Title       string
	Content     string
	IsEncrypted bool
	TimeToLive  int // minutes
}
