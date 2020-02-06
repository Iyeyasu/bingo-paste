package model

import (
	"database/sql"
	"time"

	"bingo/internal/config"
)

// User represents an authenticated user.
type User struct {
	ID           int64           `db:"id"`
	TimeCreated  time.Time       `db:"time_created"`
	UID          string          `db:"uid"`
	Name         string          `db:"name"`
	Email        sql.NullString  `db:"email"`
	PasswordHash sql.NullString  `db:"password_hash"`
	AuthMode     config.AuthMode `db:"auth_mode"`
	Role         config.Role     `db:"role"`
	Theme        config.Theme    `db:"theme"`
}

// UserTemplate represents user changes to be committed to the database.
type UserTemplate struct {
	ID          sql.NullInt64
	TimeCreated sql.NullTime
	UID         sql.NullString
	Name        sql.NullString
	Email       sql.NullString
	Password    sql.NullString
	AuthMode    sql.NullInt32
	Role        sql.NullInt32
	Theme       sql.NullInt32
}
