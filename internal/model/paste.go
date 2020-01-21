package model

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	Title             string
	Content           string
	IsEncrypted       bool
	MinutesPreserved  int
	MinutesToPreserve int // minutes
}
