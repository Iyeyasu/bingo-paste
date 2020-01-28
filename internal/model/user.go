package model

// User represents a user.
type User struct {
	// Metadata
	ID                 int64
	TimeCreatedSeconds int64

	// Settings
	Name               string
	Mail               string
	AuthenticationType string
	PrivilegeLevel     string
	Theme              string
}
