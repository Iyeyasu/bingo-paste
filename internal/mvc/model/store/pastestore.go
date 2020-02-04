package store

import (
	"database/sql"
	"fmt"
	"time"

	"bingo/internal/config"
	"bingo/internal/mvc/model"
	"bingo/internal/util/fmtutil"
	"bingo/internal/util/log"
)

var (
	deleteExpiredInterval = time.Minute * time.Duration(5)
)

// PasteStore is the store for pastes.
type PasteStore struct {
	query *PasteQuery
}

// NewPasteStore creates a new PasteStore instance.
func NewPasteStore(db *sql.DB) *PasteStore {
	log.Debug("Initializing paste store")

	store := new(PasteStore)
	store.query = NewPasteQuery(db)

	if config.Get().Expiry.Enabled {
		go store.monitorExpired()
	}

	return store
}

// Insert inserts a new paste to the database.
func (store *PasteStore) Insert(pasteTmpl *model.PasteTemplate) (*model.Paste, error) {
	log.Debug("Inserting new paste to database")

	timeCreated := time.Now().UTC()
	timeExpires := timeCreated.Add(pasteTmpl.Duration)
	formatted := fmtutil.FormatCode(pasteTmpl.Language, pasteTmpl.RawContent)
	row := store.query.insert.QueryRow(
		pasteTmpl.Title,
		pasteTmpl.RawContent,
		formatted,
		pasteTmpl.Visibility,
		timeCreated,
		sql.NullTime{Time: timeExpires, Valid: pasteTmpl.Duration > 0},
		pasteTmpl.Language,
	)

	return store.scanRow(row)
}

// FindByID returns the paste with the given id from the database.
func (store *PasteStore) FindByID(id int64) (*model.Paste, error) {
	log.Debugf("Retrieving paste %d from database", id)

	row := store.query.findByID.QueryRow(id, time.Now().UTC())
	return store.scanRow(row)
}

// FindRange returns a slice of public pastes sorted by their creation time.
func (store *PasteStore) FindRange(limit int64, offset int64) ([]*model.Paste, error) {
	log.Debugf("Retrieving %d public pastes starting from paste number %d from database", limit, offset)

	rows, err := store.query.findRange.Query(config.VisibilityListed, time.Now().UTC(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pastes: %s", err)
	}

	return store.scanRows(rows)
}

// Search returns a list of public pastes sorted by their creation time and matching given filter.
func (store *PasteStore) Search(filter string, limit int64, offset int64) ([]*model.Paste, error) {
	log.Debugf("Retrieving %d public pastes starting from paste number %d and matching matching '%s' from database", limit, offset, filter)

	rows, err := store.query.search.Query(config.VisibilityListed, time.Now().UTC(), filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pastes: %s", err)
	}

	return store.scanRows(rows)
}

// Delete deletes the paste with the given id from the database.
func (store *PasteStore) Delete(id int64) error {
	log.Debugf("Deleting paste %d from database", id)

	_, err := store.query.delete.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete paste %d: %s", id, err)
	}

	return nil
}

func (store *PasteStore) scanRows(rows *sql.Rows) ([]*model.Paste, error) {
	pastes := []*model.Paste{}
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

func (store *PasteStore) scanRow(row model.Scannable) (*model.Paste, error) {
	paste := new(model.Paste)

	var timeExpires sql.NullTime
	err := row.Scan(
		&paste.ID,
		&paste.Title,
		&paste.RawContent,
		&paste.FormattedContent,
		&paste.Visibility,
		&paste.TimeCreated,
		&timeExpires,
		&paste.Language)

	if err != nil {
		return nil, fmt.Errorf("failed to scan paste row: %s", err)
	}

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
	result, err := store.query.deleteExpired.Exec(time.Now().UTC())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired pastes: %s", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to count expired pastes: %s", err)
	}
	return count, nil
}
