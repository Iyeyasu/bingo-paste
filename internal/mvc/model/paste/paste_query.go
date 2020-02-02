package model

import (
	"database/sql"
	"fmt"

	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

var pasteColumns = "title, raw_content, formatted_content, is_public, time_created_sec, time_expires_sec, language"

// PasteQuery is the query for pastes.
type PasteQuery struct {
	findByID      *sql.Stmt
	findRange     *sql.Stmt
	search        *sql.Stmt
	insert        *sql.Stmt
	delete        *sql.Stmt
	deleteExpired *sql.Stmt
}

// NewPasteQuery creates a new PasteQuery instance.
func NewPasteQuery(db *sql.DB) *PasteQuery {
	log.Debug("Initializing paste query")

	query := new(PasteQuery)
	createTable(db)
	query.findByID = createFindByIDStatement(db)
	query.findRange = createFindRangeStatement(db)
	query.search = createSearchStatement(db)
	query.insert = createInsertStatement(db)
	query.delete = createDeleteStatement(db)
	query.deleteExpired = createDeleteExpiredStatement(db)

	return query
}

func createTable(db *sql.DB) {
	log.Debug("Creating table 'pastes'")

	q := `
		CREATE SEQUENCE IF NOT EXISTS pastes_id_seq AS bigint;

		CREATE TABLE IF NOT EXISTS pastes (
			id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('pastes_id_seq')),
			title text NOT NULL,
			raw_content text NOT NULL,
			formatted_content text NOT NULL,
			is_public bool NOT NULL,
			time_created_sec bigint NOT NULL,
			time_expires_sec bigint NOT NULL,
			language text NOT NULL,
			tsv TSVECTOR
		);

		ALTER SEQUENCE pastes_id_seq OWNED BY pastes.id;

		CREATE INDEX IF NOT EXISTS index_pastes_expires ON pastes (time_expires_sec, is_public);

		CREATE INDEX IF NOT EXISTS index_pastes_tsv ON pastes USING GIN(tsv)
		`
	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table 'pastes': %s", err)
	}
}

func createInsertStatement(db *sql.DB) *sql.Stmt {
	// We replace all periods with space as otherwise postgres won't recognize
	// period as a delimiter for full text search.
	q := fmt.Sprintf(`
		INSERT INTO pastes (%s, tsv)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
			setweight(to_tsvector($1), 'A')
			|| setweight(to_tsvector(replace($2, '.', ' ')), 'B')
			|| setweight(to_tsvector('simple', $7), 'C'))
		RETURNING id, %s
		`, pasteColumns, pasteColumns)

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to get insert paste statement: %s", err)
	}
	return stmt
}

func createFindByIDStatement(db *sql.DB) *sql.Stmt {
	q := fmt.Sprintf(`
		SELECT id, %s
		FROM pastes
		WHERE id = $1
		AND time_expires_sec > $2
		`, pasteColumns)

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to get find paste by id statement: %s", err)
	}
	return stmt
}

func createFindRangeStatement(db *sql.DB) *sql.Stmt {
	q := fmt.Sprintf(`
		SELECT id, %s
		FROM pastes
		WHERE time_expires_sec > $1
		AND is_public = TRUE
		ORDER BY time_created_sec DESC, id ASC
		LIMIT $2 OFFSET $3
		`, pasteColumns)

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to get find paste range statement: %s", err)
	}
	return stmt
}

func createSearchStatement(db *sql.DB) *sql.Stmt {
	q := fmt.Sprintf(`
		SELECT id, %s
		FROM pastes
		WHERE time_expires_sec > $1
		AND is_public = TRUE
		AND tsv @@ plainto_tsquery($2)
		ORDER BY time_created_sec DESC, id ASC
		LIMIT $3 OFFSET $4
		`, pasteColumns)

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to get search paste statement: %s", err)
	}
	return stmt
}

func createDeleteStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("DELETE FROM pastes WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to get delete paste statement: %s", err)
	}
	return stmt
}

func createDeleteExpiredStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("DELETE FROM pastes WHERE time_expires_sec <= $1")
	if err != nil {
		log.Fatalf("Failed to get delete expired pastes statement: %s", err)
	}
	return stmt
}
