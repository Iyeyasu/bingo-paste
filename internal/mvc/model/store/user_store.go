package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/mvc/model"
	"github.com/Iyeyasu/bingo-paste/internal/util/auth"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

// UserStore is the store for users.
type UserStore struct {
	query *userQuery
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *sql.DB) *UserStore {
	log.Debug("Initializing user store")
	store := new(UserStore)
	store.query = newUserQuery(db)
	log.Debugf("User store initialized (%d users)", store.Count())
	return store
}

// Count returns the number of users.
func (store *UserStore) Count() int64 {
	count := int64(0)
	err := store.query.count.QueryRow().Scan(&count)
	if err != nil {
		log.Fatalf("Failed to get user count: %s", err)
	}
	return count
}

// FindByID returns the user with the given id from the database.
func (store *UserStore) FindByID(id int64) (*model.User, error) {
	log.Debugf("Retrieving user %d from database", id)

	row := store.query.findByID.QueryRow(id)
	return store.scanRow(row)
}

// FindByName returns the user with the given name from the database.
func (store *UserStore) FindByName(name string) (*model.User, error) {
	log.Debugf("Retrieving user with name '%s' from database", name)

	row := store.query.findByName.QueryRow(name)
	return store.scanRow(row)
}

// FindByEmail returns the user with the given email from the database.
func (store *UserStore) FindByEmail(email string) (*model.User, error) {
	log.Debugf("Retrieving user with mail '%s' from database", email)

	row := store.query.findByEmail.QueryRow(email)
	return store.scanRow(row)
}

// FindRange returns a slice of public users sorted by their creation time.
func (store *UserStore) FindRange(limit int64, offset int64) ([]*model.User, error) {
	log.Debugf("Retrieving %d public users starting from user number %d from database", limit, offset)

	rows, err := store.query.findRange.Query(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %s", err)
	}

	return store.scanRows(rows)
}

// Insert inserts a new user to the database.
func (store *UserStore) Insert(userTmpl *model.UserTemplate) (*model.User, error) {
	log.Debug("Inserting new user to database")
	log.Tracef("%+v", userTmpl)

	passwordHash, err := auth.HashPassword(userTmpl.Password.String)
	if err != nil {
		return nil, err
	}

	row := store.query.insert.QueryRow(
		time.Now().Unix(),
		passwordHash,
		userTmpl.Name,
		userTmpl.Email,
		userTmpl.AuthMode,
		userTmpl.AuthExternalID,
		userTmpl.Role,
		userTmpl.Theme,
	)

	return store.scanRow(row)
}

// Update Updates an existing user in the database.
func (store *UserStore) Update(userTmpl *model.UserTemplate) (*model.User, error) {
	log.Debug("Updating existing user in the database")
	log.Tracef("%+v", userTmpl)

	var passwordHash sql.NullString
	if userTmpl.Password.Valid {
		hash, err := auth.HashPassword(userTmpl.Password.String)
		if err != nil {
			return nil, err
		}

		passwordHash.String = hash
		passwordHash.Valid = true
	}

	row := store.query.update.QueryRow(
		userTmpl.ID,
		passwordHash,
		userTmpl.Name,
		userTmpl.Email,
		userTmpl.AuthMode,
		userTmpl.AuthExternalID,
		userTmpl.Role,
		userTmpl.Theme,
	)

	return store.scanRow(row)
}

// Delete deletes the user with the given id from the database.
func (store *UserStore) Delete(id int64) error {
	log.Debugf("Deleting user %d from database", id)

	_, err := store.query.delete.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete user %d: %s", id, err)
	}

	return nil
}

func (store *UserStore) scanRows(rows *sql.Rows) ([]*model.User, error) {
	users := []*model.User{}
	defer rows.Close()
	for rows.Next() {
		user, err := store.scanRow(rows)
		if err == nil {
			users = append(users, user)
		} else {
			log.Errorf(err.Error())
		}
	}

	err := rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to scan user rows: %s", err)
	}

	return users, nil
}

func (store *UserStore) scanRow(row model.Scannable) (*model.User, error) {
	user := new(model.User)
	timeCreated := int64(0)
	authExternalID := []byte{}
	err := row.Scan(
		&user.ID,
		&timeCreated,
		&user.PasswordHash,
		&user.Name,
		&user.Email,
		&user.AuthMode,
		&authExternalID,
		&user.Role,
		&user.Theme,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %s", err)
	}

	user.TimeCreated = time.Unix(timeCreated, 0)
	user.AuthExternalID = string(authExternalID)
	log.Tracef("Retrieved user %+v", user)
	return user, nil
}
