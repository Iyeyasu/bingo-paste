package model

import (
	"database/sql"
	"log"
	"time"
)

// PasteStore is the store for pastes.
type PasteStore struct {
	selectStmt *sql.Stmt
	insertStmt *sql.Stmt
}

// NewStore creates a new PasteStore instance.
func NewStore(db *sql.DB) *PasteStore {
	log.Println("New paste store")

	createPseudoEncrypt(db)
	createTable(db)

	store := new(PasteStore)
	store.selectStmt = getSelectStatement(db)
	store.insertStmt = getInsertStatement(db)
	return store
}

// Insert inserts a new paste to the database.
func (store *PasteStore) Insert(paste *Paste) (int64, error) {
	log.Printf("INSERT paste %+v", paste)

	var id int64
	err := store.insertStmt.QueryRow(
		paste.Title,
		paste.Content,
		paste.IsPublic,
		time.Now().Unix(),
		paste.LifetimeSeconds).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Select returns the paste from the database with the given id.
func (store *PasteStore) Select(id int64) (*Paste, error) {
	log.Printf("SELECT paste %d", id)

	paste := NewPaste()
	err := store.selectStmt.QueryRow(id).Scan(
		&paste.ID,
		&paste.Title,
		&paste.Content,
		&paste.IsPublic,
		&paste.TimeCreatedSeconds,
		&paste.LifetimeSeconds)

	if err != nil {
		return nil, err
	}

	return paste, nil
}

func createTable(db *sql.DB) {
	log.Printf("Creating table")

	q := "CREATE SEQUENCE IF NOT EXISTS pastes_id_seq AS bigint"
	_, err := db.Exec(q)
	if err != nil {
		log.Fatalln(err)
	}

	q = `
	CREATE TABLE IF NOT EXISTS pastes (
		id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('pastes_id_seq')),
		title text NOT NULL,
		content text NOT NULL,
		is_public bool,
		time_created_seconds bigint,
		lifetime_seconds bigint)
	`
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalln(err)
	}

	q = "ALTER SEQUENCE pastes_id_seq OWNED BY pastes.id"
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalln(err)
	}
}

// https://stackoverflow.com/questions/12761346/pseudo-encrypt-function-in-plpgsql-that-takes-bigint/12761795#12761795
// Creates a function that maps big integers to another seemingly random big integer.
// Used to make sure the ids of pastes are seemingly random.
func createPseudoEncrypt(db *sql.DB) {
	log.Printf("Creating pseudo encrypt function")

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
		log.Fatalln(err)
	}
}

func getInsertStatement(db *sql.DB) *sql.Stmt {
	log.Printf("Getting prepared insert statement")

	query := "INSERT INTO pastes (title, content, is_public, time_created_seconds, lifetime_seconds) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func getSelectStatement(db *sql.DB) *sql.Stmt {
	log.Printf("Getting prepared select statement")

	query := "SELECT id, title, content, is_public, time_created_seconds, lifetime_seconds FROM pastes WHERE id = $1"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
