package model

import (
	"database/sql"
	"time"
)

// Role represents the access level of the user.
type Role int

// AuthType represents the access level of the user.
type AuthType int

// Theme is the theme used by the user.
type Theme int

const (
	// RoleAdmin has the power to change site configuration and users.
	RoleAdmin Role = iota

	// RoleEditor has the power to create new pastes.
	RoleEditor = iota

	// RoleViewer can view pastes.
	RoleViewer = iota
)

const (
	// AuthStandard users log in using a stored password.
	AuthStandard AuthType = iota

	// AuthLDAP users log in using their LDAP credentials.
	AuthLDAP = iota
)

const (
	// ThemeDark represents a dark GUI theme.
	ThemeDark Theme = iota

	// ThemeLight represents a light GUI theme.
	ThemeLight = iota
)

// User represents an authenticated user.
type User struct {
	ID             int64
	TimeCreated    time.Time
	PasswordHash   []byte
	Name           string
	Mail           string
	Role           Role
	Theme          Theme
	AuthType       AuthType
	AuthExternalID string
}

// UserModel represents user changes to be committed to the database.
type UserModel struct {
	ID             sql.NullInt64
	Password       sql.NullString
	Name           sql.NullString
	Mail           sql.NullString
	AuthType       sql.NullInt32
	AuthExternalID sql.NullString
	Role           sql.NullInt32
	Theme          sql.NullInt32
}
