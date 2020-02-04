package session

import (
	"database/sql"
	"log"
	"time"

	"github.com/alexedwards/scs/postgresstore"
)

// PostgresStore is a wrapper around postgresstore.PostgresStore that creates
// the required sessions table.
type PostgresStore struct {
	store *postgresstore.PostgresStore
}

// NewPostgresStore creates a new PostgresStore.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	store := new(PostgresStore)
	store.store = postgresstore.New(db)
	createTable(db)
	return store
}

// Find returns the data for a given session token from the PostgresStore instance.
// If the session token is not found or is expired, the returned exists flag
// will be set to false.
func (store *PostgresStore) Find(token string) (b []byte, exists bool, err error) {
	return store.store.Find(token)
}

// Commit adds a session token and data to the PostgresStore instance with the
// given expiry time. If the session token already exists then the data and
// expiry time are updated.
func (store *PostgresStore) Commit(token string, b []byte, expiry time.Time) error {
	return store.store.Commit(token, b, expiry)
}

// Delete removes a session token and corresponding data from the PostgresStore
// instance.
func (store *PostgresStore) Delete(token string) error {
	return store.store.Delete(token)
}

func createTable(db *sql.DB) {
	q := `
	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		data BYTEA NOT NULL,
		expiry TIMESTAMPTZ NOT NULL
	);
	CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
	`

	_, err := db.Exec(q)
	if err != nil {
		log.Fatalln("Failed to create table 'sessions':", err)
	}
}
