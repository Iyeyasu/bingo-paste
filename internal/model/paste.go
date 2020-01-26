package model

import (
	"encoding/json"
	"time"
)

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

// MarshalBinary converts the paste to a binary array.
func (paste *Paste) MarshalBinary() ([]byte, error) {
	return json.Marshal(paste)
}

// UnmarshalBinary converts a binary array to paste.
func (paste *Paste) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, paste)
}
