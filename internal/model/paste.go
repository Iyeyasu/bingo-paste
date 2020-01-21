package model

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	Content            string
	IsPublic           bool
	TimeCreatedSeconds int64
	LifetimeSeconds    int64
}
