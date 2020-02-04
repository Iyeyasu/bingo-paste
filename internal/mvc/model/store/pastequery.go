package store

import (
	"database/sql"
	"fmt"

	"bingo/internal/util/log"
)

var pasteColumns = "title, raw_content, formatted_content, visibility, time_created, time_expires, language"

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
	query.createTable(db)
	query.findByID = query.createFindByIDStatement(db)
	query.findRange = query.createFindRangeStatement(db)
	query.search = query.createSearchStatement(db)
	query.insert = query.createInsertStatement(db)
	query.delete = query.createDeleteStatement(db)
	query.deleteExpired = query.createDeleteExpiredStatement(db)

	return query
}

func (q *PasteQuery) createTable(db *sql.DB) {
	log.Debug("Creating table 'pastes'")

	query := `
		CREATE SEQUENCE IF NOT EXISTS pastes_id_seq AS bigint;

		CREATE TABLE IF NOT EXISTS pastes (
			id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('pastes_id_seq')),
			title text NOT NULL,
			raw_content text NOT NULL,
			formatted_content text NOT NULL,
			visibility int NOT NULL,
			time_created timestamptz NOT NULL,
			time_expires timestamptz,
			language text NOT NULL,
			tsv TSVECTOR
		);

		ALTER SEQUENCE pastes_id_seq OWNED BY pastes.id;
		CREATE INDEX IF NOT EXISTS pastes_time_expires_id_idx ON pastes (time_expires, id);
		CREATE INDEX IF NOT EXISTS pastes_tsv_idx ON pastes USING GIN(tsv)
		`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table 'pastes': %s", err)
	}
}

func (q *PasteQuery) createInsertStatement(db *sql.DB) *sql.Stmt {
	// We replace all periods with space as otherwise postgres won't recognize
	// period as a delimiter for full text search.
	query := fmt.Sprintf(`
		INSERT INTO pastes (%s, tsv)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
			setweight(to_tsvector($1), 'A')
			|| setweight(to_tsvector(replace($2, '.', ' ')), 'B')
			|| setweight(to_tsvector('simple', $7), 'C'))
		RETURNING id, %s
		`, pasteColumns, pasteColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get insert paste statement: %s", err)
	}
	return stmt
}

func (q *PasteQuery) createFindByIDStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
		SELECT id, %s
		FROM pastes
		WHERE id = $1
		AND (time_expires IS NULL OR time_expires > $2)
		`, pasteColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get find paste by id statement: %s", err)
	}
	return stmt
}

func (q *PasteQuery) createFindRangeStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
		SELECT id, %s
		FROM pastes
		WHERE visibility >= $1
		AND (time_expires IS NULL OR time_expires > $2)
		ORDER BY time_created DESC, id ASC
		LIMIT $3 OFFSET $4
		`, pasteColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get find paste range statement: %s", err)
	}
	return stmt
}

func (q *PasteQuery) createSearchStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
		SELECT id, %s
		FROM pastes
		WHERE visibility >= $1
		AND (time_expires IS NULL OR time_expires > $2)
		AND tsv @@ plainto_tsquery($3)
		ORDER BY time_created DESC, id ASC
		LIMIT $4 OFFSET $5
		`, pasteColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get search paste statement: %s", err)
	}
	return stmt
}

func (q *PasteQuery) createDeleteStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("DELETE FROM pastes WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to get delete paste statement: %s", err)
	}
	return stmt
}

func (q *PasteQuery) createDeleteExpiredStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("DELETE FROM pastes WHERE time_expires <= $1")
	if err != nil {
		log.Fatalf("Failed to get delete expired pastes statement: %s", err)
	}
	return stmt
}
