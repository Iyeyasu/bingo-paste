package model

import (
	"time"

	"bingo/internal/config"
)

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID               int64
	TimeCreated      time.Time
	Title            string
	RawContent       string
	FormattedContent string
	Visibility       config.Visibility
	Language         string
}

// PasteTemplate represents paste changes to be committed to the database.
type PasteTemplate struct {
	Title      string
	RawContent string
	Visibility config.Visibility
	Language   string
	Duration   time.Duration
}
