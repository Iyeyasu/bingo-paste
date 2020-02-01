package model

import (
	"database/sql"
	"fmt"
	"time"

	util "github.com/Iyeyasu/bingo-paste/internal/util/html"
	log "github.com/sirupsen/logrus"
)

var (
	deleteExpiredInterval = time.Minute * time.Duration(5)
)

// PasteStore is the store for pastes.
type PasteStore struct {
	selectStmt        *sql.Stmt
	selectListStmt    *sql.Stmt
	searchStmt        *sql.Stmt
	insertStmt        *sql.Stmt
	deleteStmt        *sql.Stmt
	deleteExpiredStmt *sql.Stmt
}

// NewPasteStore creates a new PasteStore instance.
func NewPasteStore(db *sql.DB) *PasteStore {
	log.Debug("Initializing paste store")

	store := new(PasteStore)
	store.createTable(db)
	store.createInsertStatement(db)
	store.createSelectStatement(db)
	store.createListStatement(db)
	store.createSearchStatement(db)
	store.createDeleteStatement(db)
	store.createDeleteExpiredStatement(db)

	go store.monitorExpired()

	return store
}

// Insert inserts a new paste to the database.
func (store *PasteStore) Insert(paste *Paste) (*Paste, error) {
	log.Debug("Inserting new paste to database")

	timeCreated := time.Now().Unix()
	timeExpires := timeCreated + int64(paste.Duration.Seconds())
	row := store.insertStmt.QueryRow(
		paste.Title,
		paste.RawContent,
		util.HighlightSyntax(paste.Language, paste.RawContent),
		paste.IsPublic,
		timeCreated,
		timeExpires,
		paste.Language,
	)

	return store.scanRow(row)
}

// Get returns the paste with the given id from the database.
func (store *PasteStore) Get(id int64) (*Paste, error) {
	log.Debugf("Retrieving paste %d from database", id)

	row := store.selectStmt.QueryRow(id, time.Now().Unix())
	return store.scanRow(row)
}

// GetList returns a slice of public pastes sorted by their creation time.
func (store *PasteStore) GetList(limit int64, offset int64) ([]*Paste, error) {
	log.Debugf("Retrieving %d public pastes starting from paste number %d from database", limit, offset)

	rows, err := store.selectListStmt.Query(time.Now().Unix(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pastes: %s", err)
	}

	return store.scanRows(rows)
}

// Search returns a list of public pastes sorted by their creation time and matching given filter.
func (store *PasteStore) Search(filter string, limit int64, offset int64) ([]*Paste, error) {
	log.Debugf("Retrieving %d public pastes starting from paste number %d and matching matching '%s' from database", limit, offset, filter)

	rows, err := store.searchStmt.Query(time.Now().Unix(), filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pastes: %s", err)
	}

	return store.scanRows(rows)
}

// Delete deletes the paste with the given id from the database.
func (store *PasteStore) Delete(id int64) error {
	log.Debugf("Deleting paste %d from database", id)

	_, err := store.deleteStmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete paste %d: %s", id, err)
	}

	return nil
}

func (store *PasteStore) scanRows(rows *sql.Rows) ([]*Paste, error) {
	pastes := []*Paste{}
	defer rows.Close()
	for rows.Next() {
		paste, err := store.scanRow(rows)
		if err == nil {
			pastes = append(pastes, paste)
		} else {
			log.Errorf(err.Error())
		}
	}

	err := rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to scan paste rows: %s", err)
	}

	return pastes, nil
}

func (store *PasteStore) scanRow(row Scannable) (*Paste, error) {
	paste := new(Paste)
	timeCreated := int64(0)
	timeExpires := int64(0)
	err := row.Scan(
		&paste.ID,
		&paste.Title,
		&paste.RawContent,
		&paste.FormattedContent,
		&paste.IsPublic,
		&timeCreated,
		&timeExpires,
		&paste.Language)

	if err != nil {
		return nil, fmt.Errorf("failed to scan paste row: %s", err)
	}

	paste.TimeCreated = time.Unix(timeCreated, 0)
	paste.Duration = time.Second * time.Duration(timeExpires-time.Now().Unix())
	return paste, nil
}

func (store *PasteStore) monitorExpired() {
	log.Debug("Monitoring expired pastes")

	for range time.Tick(deleteExpiredInterval) {
		count, err := store.deleteExpired()
		if err != nil {
			log.Errorf(err.Error())
		} else {
			log.Debugf("Deleted %d expired pastes", count)
		}
	}
}

func (store *PasteStore) deleteExpired() (int64, error) {
	result, err := store.deleteExpiredStmt.Exec(time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired pastes: %s", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to count expired pastes: %s", err)
	}
	return count, nil
}

func (store *PasteStore) createTable(db *sql.DB) {
	log.Debug("Creating table 'pastes'")

	q := "CREATE SEQUENCE IF NOT EXISTS pastes_id_seq AS bigint"
	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create sequence 'pastes_id_seq': %s", err)
	}

	q = `
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
	)
	`
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table 'pastes': %s", err)
	}

	q = "ALTER SEQUENCE pastes_id_seq OWNED BY pastes.id"
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to assign sequence 'pastes_id_seq': %s", err)
	}

	q = "CREATE INDEX IF NOT EXISTS index_pastes_expires ON pastes (time_expires_sec, is_public)"
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create index 'index_pastes_expires': %s", err)
	}

	q = "CREATE INDEX IF NOT EXISTS index_pastes_tsv ON pastes USING GIN(tsv);"
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create index 'index_pastes_tsv': %s", err)
	}
}

func (store *PasteStore) createInsertStatement(db *sql.DB) {
	// We replace all periods with space as otherwise postgres won't recognize
	// period as a delimiter for full text search.
	query := `
	INSERT INTO pastes
	(title, raw_content, formatted_content, is_public, time_created_sec, time_expires_sec, language, tsv)
	VALUES ($1, $2, $3, $4, $5, $6, $7,
		setweight(to_tsvector($1), 'A')
		|| setweight(to_tsvector(replace($2, '.', ' ')), 'B')
		|| setweight(to_tsvector('simple', $7), 'C'))
	RETURNING id, title, raw_content, formatted_content, is_public, time_created_sec, time_expires_sec, language
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get insert statement: %s", err)
	}
	store.insertStmt = stmt
}

func (store *PasteStore) createSelectStatement(db *sql.DB) {
	query := `
	SELECT id, title, raw_content, formatted_content, is_public, time_created_sec, time_expires_sec, language
	FROM pastes
	WHERE id = $1
	AND time_expires_sec> $2
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get select statement: %s", err)
	}
	store.selectStmt = stmt
}

func (store *PasteStore) createListStatement(db *sql.DB) {
	query := `
	SELECT id, title, raw_content, formatted_content, is_public, time_created_sec, time_expires_sec, language
	FROM pastes
	WHERE time_expires_sec> $1
	AND is_public = TRUE
	ORDER BY time_created_sec DESC, id ASC
	LIMIT $2 OFFSET $3
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get select list statement: %s", err)
	}
	store.selectListStmt = stmt
}

func (store *PasteStore) createSearchStatement(db *sql.DB) {
	query := `
	SELECT id, title, raw_content, formatted_content, is_public, time_created_sec, time_expires_sec, language
	FROM pastes
	WHERE time_expires_sec> $1
	AND is_public = TRUE
	AND tsv @@ plainto_tsquery($2)
	ORDER BY time_created_sec DESC, id ASC
	LIMIT $3 OFFSET $4
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get select search statement: %s", err)
	}
	store.searchStmt = stmt
}

func (store *PasteStore) createDeleteStatement(db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM pastes WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to get delete statement: %s", err)
	}
	store.deleteStmt = stmt
}

func (store *PasteStore) createDeleteExpiredStatement(db *sql.DB) {
	stmt, err := db.Prepare("DELETE FROM pastes WHERE time_expires_sec < $1")
	if err != nil {
		log.Fatalf("Failed to get delete expired statement: %s", err)
	}
	store.deleteExpiredStmt = stmt
}
