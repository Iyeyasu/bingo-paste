package model

// Scannable is a common interface for sql.Row and sql.Rows.
type Scannable interface {
	Scan(dest ...interface{}) error
}
