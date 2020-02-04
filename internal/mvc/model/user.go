package model

import (
	"database/sql"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"
)

// User represents an authenticated user.
type User struct {
	ID             int64
	TimeCreated    time.Time
	PasswordHash   string
	Name           string
	Email          string
	Role           config.Role
	Theme          config.Theme
	AuthMode       config.AuthMode
	AuthExternalID string
}

// UserTemplate represents user changes to be committed to the database.
type UserTemplate struct {
	ID             sql.NullInt64
	TimeCreated    sql.NullTime
	Password       sql.NullString
	Name           sql.NullString
	Email          sql.NullString
	AuthMode       sql.NullInt32
	AuthExternalID sql.NullString
	Role           sql.NullInt32
	Theme          sql.NullInt32
}
