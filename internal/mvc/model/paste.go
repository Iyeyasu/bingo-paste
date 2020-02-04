package model

import (
	"database/sql"
	"time"

	"bingo/internal/config"
)

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID               int64             `db:"id"`
	TimeCreated      time.Time         `db:"time_created"`
	Title            string            `db:"title"`
	RawContent       string            `db:"raw_content"`
	FormattedContent string            `db:"formatted_content"`
	Language         string            `db:"language"`
	TimeExpires      sql.NullTime      `db:"time_expires"`
	Visibility       config.Visibility `db:"visibility"`
}

// PasteTemplate represents paste changes to be committed to the database.
type PasteTemplate struct {
	Title      string
	RawContent string
	Visibility config.Visibility
	Language   string
	Duration   time.Duration
}
