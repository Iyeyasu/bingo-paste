package model

import "time"

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID               int64
	Title            string
	RawContent       string
	FormattedContent string
	IsPublic         bool
	Language         string
	TimeCreated      time.Time
	Duration         time.Duration
}
