package model

import (
	"database/sql"
	"log"
	"time"
)

// PasteDB is the interface for the stored pastes.
type PasteDB struct {
	db        *sql.DB
	statement *sql.Stmt
}

// Initialize initializes the paste database.
func (pastes *PasteDB) Initialize(db *sql.DB) {
	pastes.createTable()
	pastes.db = db
	pastes.statement = pastes.getPreparedStatement()
}

func (pastes *PasteDB) getPreparedStatement() *sql.Stmt {
	query := "SELECT content, is_public, time_created_seconds, lifetime_minutes FROM pastes WHERE id = $1"
	stmt, err := pastes.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func (pastes *PasteDB) createTable() {
	q := "CREATE TABLE IF NOT EXISTS pastes (id integer PRIMARY KEY, content text NOT NULL, is_public bool, time_created_seconds bigint, lifetime_seconds bigint)"
	_, err := pastes.db.Exec(q)
	if err != nil {
		log.Fatalln(err)
	}
}

func (pastes *PasteDB) insert(paste Paste) {
	id := uint64(0)
	_, err := pastes.db.Exec(
		"INSERT INTO pastes (id, content, is_public, time_created_seconds, lifetime_seconds) VALUES ($1, $2, $3, $4 $5)",
		id,
		paste.Content,
		paste.IsPublic,
		time.Now().Unix(),
		paste.LifetimeSeconds)

	if err != nil {
		log.Fatal(err)
	}
}

func (pastes *PasteDB) get(id uint64) Paste {
	var paste Paste

	err := pastes.statement.QueryRow(id).Scan(
		&paste.Content,
		&paste.IsPublic,
		&paste.TimeCreatedSeconds,
		&paste.LifetimeSeconds)

	if err != nil {
		log.Fatal(err)
	}

	return paste
}
