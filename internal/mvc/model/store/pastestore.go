package store

import (
	"database/sql"
	"time"

	"bingo/internal/config"
	"bingo/internal/mvc/model"
	"bingo/internal/util/fmtutil"
	"bingo/internal/util/log"

	"github.com/jmoiron/sqlx"
)

var (
	deleteExpiredInterval = time.Minute * time.Duration(5)
)

// PasteStore is the store for pastes.
type PasteStore struct {
	Database *sqlx.DB
}

// NewPasteStore creates a new PasteStore instance.
func NewPasteStore(db *sqlx.DB) *PasteStore {
	log.Debug("Initializing paste store")

	store := new(PasteStore)
	store.Database = db
	store.createTable()

	if config.Get().Expiry.Enabled {
		go store.monitorExpired()
	}

	return store
}

// Count returns the number of pastes.
func (store *PasteStore) Count() int64 {
	log.Debugf("Counting number of pastes")

	var count int64
	store.Database.Get(&count, "SELECT COUNT(*) FROM pastes")
	return count
}

// FindByID returns the paste with the given id from the database.
func (store *PasteStore) FindByID(id int64) (*model.Paste, error) {
	log.Debugf("Retrieving paste %d from database", id)

	query := `
		SELECT id, time_created, title, raw_content, formatted_content, language, time_expires, visibility
		FROM pastes
		WHERE id = $1
		AND (time_expires IS NULL OR time_expires > $2)
		`

	paste := new(model.Paste)
	err := store.Database.Get(paste, query, id, time.Now().UTC())
	return paste, err
}

// FindRange returns a slice of public pastes sorted by their creation time.
func (store *PasteStore) FindRange(limit int64, offset int64) ([]model.Paste, error) {
	log.Debugf("Retrieving %d public pastes starting from paste number %d from database", limit, offset)

	query := `
		SELECT id, time_created, title, raw_content, formatted_content, language, time_expires, visibility
		FROM pastes
		WHERE visibility >= $1
		AND (time_expires IS NULL OR time_expires > $2)
		ORDER BY time_created DESC, id ASC
		LIMIT $3 OFFSET $4
		`

	pastes := []model.Paste{}
	err := store.Database.Select(&pastes, query, config.VisibilityListed, time.Now().UTC(), limit, offset)
	return pastes, err
}

// Search returns a list of public pastes sorted by their creation time and matching given filter.
func (store *PasteStore) Search(filter string, limit int64, offset int64) ([]model.Paste, error) {
	log.Debugf("Retrieving %d public pastes starting from paste number %d and matching matching '%s' from database", limit, offset, filter)

	query := `
		SELECT id, time_created, title, raw_content, formatted_content, language, time_expires, visibility
		FROM pastes
		WHERE visibility >= $1
		AND (time_expires IS NULL OR time_expires > $2)
		AND tsv @@ plainto_tsquery($3)
		ORDER BY time_created DESC, id ASC
		LIMIT $4 OFFSET $5
		`

	pastes := []model.Paste{}
	err := store.Database.Select(&pastes, query, config.VisibilityListed, time.Now().UTC(), filter, limit, offset)
	return pastes, err
}

// Delete deletes the paste with the given id from the database.
func (store *PasteStore) Delete(id int64) error {
	log.Debugf("Deleting paste %d from database", id)

	_, err := store.Database.Exec("DELETE FROM pastes WHERE id = $1", id)
	return err
}

// Insert inserts a new paste to the database.
func (store *PasteStore) Insert(pasteTmpl *model.PasteTemplate) (*model.Paste, error) {
	log.Debug("Inserting new paste to database")

	query := `
		INSERT INTO pastes (time_created, title, raw_content, formatted_content, language, time_expires, visibility, tsv)
		VALUES ($1, $2, $3, $4, $5, $6, $7,
			setweight(to_tsvector($2), 'A')
			|| setweight(to_tsvector(replace($3, '.', ' ')), 'B')
			|| setweight(to_tsvector('simple', $5), 'C'))
		RETURNING id, time_created, title, raw_content, formatted_content, language, time_expires, visibility
		`

	paste := new(model.Paste)
	timeCreated := time.Now().UTC()
	timeExpires := timeCreated.Add(pasteTmpl.Duration)
	formatted := fmtutil.FormatCode(pasteTmpl.Language, pasteTmpl.RawContent)
	err := store.Database.QueryRowx(
		query,
		timeCreated,
		pasteTmpl.Title,
		pasteTmpl.RawContent,
		formatted,
		pasteTmpl.Language,
		sql.NullTime{Time: timeExpires, Valid: pasteTmpl.Duration > 0},
		pasteTmpl.Visibility,
	).StructScan(paste)

	return paste, err
}

func (store *PasteStore) createTable() {
	log.Debug("Creating table 'pastes'")

	query := `
		CREATE SEQUENCE IF NOT EXISTS pastes_id_seq AS bigint;

		CREATE TABLE IF NOT EXISTS pastes (
			id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('pastes_id_seq')),
			time_created timestamptz NOT NULL,
			title text NOT NULL,
			raw_content text NOT NULL,
			formatted_content text NOT NULL,
			language text NOT NULL,
			time_expires timestamptz,
			visibility int NOT NULL,
			tsv TSVECTOR
		);

		ALTER SEQUENCE pastes_id_seq OWNED BY pastes.id;
		CREATE INDEX IF NOT EXISTS pastes_time_expires_id_idx ON pastes (time_expires, id);
		CREATE INDEX IF NOT EXISTS pastes_tsv_idx ON pastes USING GIN(tsv)
		`
	_, err := store.Database.Exec(query)
	if err != nil {
		log.Fatalln("Failed to create table 'pastes':", err)
	}
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
	result, err := store.Database.Exec("DELETE FROM pastes WHERE time_expires <= $1", time.Now().UTC())
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}
