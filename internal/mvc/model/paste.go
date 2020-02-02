package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"
)

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID               int64         `db:"id"`
	TimeCreated      time.Time     `db:"time_created_sec"`
	Title            string        `db:"title"`
	RawContent       string        `db:"raw_content"`
	FormattedContent string        `db:"formatted_content"`
	IsPublic         bool          `db:"is_public"`
	Language         string        `db:"language"`
	Duration         time.Duration `db:"duration"`
}

// MarshalBinary converts the paste to a binary array.
func (paste *Paste) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(paste); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary converts a binary array to paste.
func (paste *Paste) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	if err := gob.NewDecoder(reader).Decode(paste); err != nil {
		return err
	}
	return json.Unmarshal(data, paste)
}
