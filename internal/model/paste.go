package model

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"time"
)

// Paste represents the paste contents and surrounding metadata.
type Paste struct {
	ID               int64
	TimeCreated      time.Time
	Title            string
	RawContent       string
	FormattedContent string
	IsPublic         bool
	Language         string
	Duration         time.Duration
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
