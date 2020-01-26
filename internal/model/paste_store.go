package model

import (
	"database/sql"
	"fmt"
	"time"

	util "github.com/Iyeyasu/bingo-paste/internal/util/html"
	log "github.com/sirupsen/logrus"
)

// PasteStore is the store for pastes.
type PasteStore struct {
	selectStmt        *sql.Stmt
	insertStmt        *sql.Stmt
	deleteStmt        *sql.Stmt
	deleteExpiredStmt *sql.Stmt
}

// NewPasteStore creates a new PasteStore instance.
func NewPasteStore(db *sql.DB) *PasteStore {
	log.Debug("Initializing paste store")

	createPseudoEncrypt(db)
	createTable(db)

	store := new(PasteStore)
	store.selectStmt = getSelectStatement(db)
	store.insertStmt = getInsertStatement(db)
	store.deleteStmt = getDeleteStatement(db)
	store.deleteExpiredStmt = getDeleteExpiredStatement(db)

	go store.monitorExpired()

	return store
}

// Insert inserts a new paste to the database.
func (store *PasteStore) Insert(paste *Paste) (int64, error) {
	log.Debug("Inserting new paste")

	id := int64(0)
	err := store.insertStmt.QueryRow(
		paste.Title,
		paste.RawContent,
		util.HighlightSyntax(paste.Language, paste.RawContent),
		paste.IsPublic,
		time.Now().Unix(),
		paste.Duration,
		paste.Language).Scan(&id)

	if err != nil {
		log.Errorf("Failed to create paste: %s", err)
		return 0, err
	}

	return id, nil
}

// Select returns the paste with the given id from the database.
func (store *PasteStore) Select(id int64) (*Paste, error) {
	log.Debugf("Retrieving paste %d", id)

	paste := new(Paste)
	timeCreated := int64(0)
	err := store.selectStmt.QueryRow(id).Scan(
		&paste.ID,
		&paste.Title,
		&paste.RawContent,
		&paste.FormattedContent,
		&paste.IsPublic,
		&timeCreated,
		&paste.Duration,
		&paste.Language)
	paste.TimeCreated = time.Unix(timeCreated, 0)

	if err != nil {
		log.Debugf("Failed to retrieve paste %d: %s", id, err)
		return nil, err
	}

	return store.filterExpired(paste)
}

// Delete deletes the paste with the given id from the database.
func (store *PasteStore) Delete(id int64) error {
	log.Debugf("Deleting paste %d", id)

	_, err := store.deleteStmt.Exec(id)
	if err != nil {
		log.Debugf("Failed to delete paste %d: %s", id, err)
	}

	return err
}

func (store *PasteStore) monitorExpired() {
	log.Info("Monitoring expired pastes")

	store.deleteExpired()
	for range time.Tick(time.Hour) {
		store.deleteExpired()
	}
}

func (store *PasteStore) deleteExpired() {
	result, err := store.deleteExpiredStmt.Exec(time.Now().Unix())
	if err != nil {
		log.Errorf("Failed to delete expired pastes: %s", err)
		return
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Errorf("Failed to count expired pastes: %s", err)
		return
	}

	log.Infof("Deleted %d expired pastes", count)
}

func (store *PasteStore) filterExpired(paste *Paste) (*Paste, error) {
	log.Debugf("Checking if paste %d with life time %d has expired", paste.ID, int64(paste.Duration.Seconds()))

	timeLeft := paste.TimeCreated.Add(paste.Duration).Sub(time.Now())
	if timeLeft > 0 {
		log.Debugf("Paste %d has not expired (%s left)", paste.ID, timeLeft)
		return paste, nil
	}

	err := store.Delete(paste.ID)
	if err != nil {
		log.Debugf("Failed to delete expired paste %d: %s", paste.ID, err)
		return nil, err
	}

	if paste.Duration > 0 {
		log.Debugf("Paste %d has expired", paste.ID)
		return nil, fmt.Errorf("paste %d expired", paste.ID)
	}

	log.Debugf("Paste %d burned on read", paste.ID)
	return paste, nil
}

func createTable(db *sql.DB) {
	log.Debug("Creating table 'pastes'")

	q := "CREATE SEQUENCE IF NOT EXISTS pastes_id_seq AS bigint"
	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	q = `
	CREATE TABLE IF NOT EXISTS pastes (
		id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('pastes_id_seq')),
		title text NOT NULL,
		raw_content text NOT NULL,
		formatted_content text NOT NULL,
		language text NOT NULL,
		is_public bool,
		time_created_seconds bigint,
		duration bigint)
	`
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}

	q = "ALTER SEQUENCE pastes_id_seq OWNED BY pastes.id"
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}
}

// https://stackoverflow.com/questions/12761346/pseudo-encrypt-function-in-plpgsql-that-takes-bigint/12761795#12761795
// Creates a function that maps big integers to another seemingly random big integer.
// Used to make sure the ids of pastes are seemingly random.
func createPseudoEncrypt(db *sql.DB) {
	log.Debug("Creating pseudo encrypt function")

	q := `
	CREATE OR REPLACE FUNCTION pseudo_encrypt(VALUE bigint) returns bigint AS $$
	DECLARE
	l1 bigint;
	l2 bigint;
	r1 bigint;
	r2 bigint;
	i int:=0;
	BEGIN
		l1:= (VALUE >> 32) & 4294967295::bigint;
		r1:= VALUE & 4294967295;
		WHILE i < 3 LOOP
			l2 := r1;
			r2 := l1 # ((((1366.0 * r1 + 150889) % 714025) / 714025.0) * 32767*32767)::int;
			l1 := l2;
			r1 := r2;
			i := i + 1;
		END LOOP;
	RETURN ((l1::bigint << 32) + r1);
	END;
	$$ LANGUAGE plpgsql strict immutable;
	`
	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create function: %s", err)
	}
}

func getInsertStatement(db *sql.DB) *sql.Stmt {
	log.Trace("Getting prepared insert statement for pastes")

	query := "INSERT INTO pastes (title, raw_content, formatted_content, is_public, time_created_seconds, duration, language) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get statement: %s", err)
	}
	return stmt
}

func getSelectStatement(db *sql.DB) *sql.Stmt {
	log.Trace("Getting prepared select statement for pastes")

	query := "SELECT id, title, raw_content, formatted_content, is_public, time_created_seconds, duration, language FROM pastes WHERE id = $1"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to get statement: %s", err)
	}
	return stmt
}

func getDeleteStatement(db *sql.DB) *sql.Stmt {
	log.Trace("Getting prepared delete statement for pastes")

	stmt, err := db.Prepare("DELETE FROM pastes WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to get statement: %s", err)
	}

	return stmt
}

func getDeleteExpiredStatement(db *sql.DB) *sql.Stmt {
	log.Trace("Getting prepared delete expired statement for pastes")

	stmt, err := db.Prepare("DELETE FROM pastes WHERE duration + time_created_seconds < $1")
	if err != nil {
		log.Fatalf("Failed to get statement: %s", err)
	}

	return stmt
}
