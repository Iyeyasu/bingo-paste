package model

import (
	"database/sql"
	"fmt"
	"time"

	util "github.com/Iyeyasu/bingo-paste/internal/util/auth"
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

	if store.Count() == 0 {
		store.createInitialUser()
	}

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
func (store *UserStore) FindByID(id int64) (*User, error) {
	log.Debugf("Retrieving user %d from database", id)

	row := store.query.findByID.QueryRow(id)
	return store.scanRow(row)
}

// FindByName returns the user with the given name from the database.
func (store *UserStore) FindByName(name string) (*User, error) {
	log.Debugf("Retrieving user %d from database", name)

	row := store.query.findByName.QueryRow(name)
	return store.scanRow(row)
}

// FindByEmail returns the user with the given email from the database.
func (store *UserStore) FindByEmail(email string) (*User, error) {
	log.Debugf("Retrieving user %d from database", email)

	row := store.query.findByEmail.QueryRow(email)
	return store.scanRow(row)
}

// FindRange returns a slice of public users sorted by their creation time.
func (store *UserStore) FindRange(limit int64, offset int64) ([]*User, error) {
	log.Debugf("Retrieving %d public users starting from user number %d from database", limit, offset)

	rows, err := store.query.findRange.Query(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %s", err)
	}

	return store.scanRows(rows)
}

// Insert inserts a new user to the database.
func (store *UserStore) Insert(userTmpl *UserModel) (*User, error) {
	log.Debug("Inserting new user to database")
	log.Tracef("%+v", userTmpl)

	passwordHash, err := util.HashPassword(userTmpl.Password.String)
	if err != nil {
		return nil, err
	}

	timeCreated := time.Now().Unix()
	row := store.query.insert.QueryRow(
		timeCreated,
		passwordHash,
		userTmpl.Name,
		userTmpl.Email,
		userTmpl.AuthType,
		userTmpl.AuthExternalID,
		userTmpl.Role,
		userTmpl.Theme,
	)

	return store.scanRow(row)
}

// Update Updates an existing user in the database.
func (store *UserStore) Update(userTmpl *UserModel) (*User, error) {
	log.Debug("Updating existing user in the database")
	log.Tracef("%+v", userTmpl)

	row := store.query.update.QueryRow(
		userTmpl.ID,
		nil,
		userTmpl.Name,
		userTmpl.Email,
		userTmpl.AuthType,
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

func (store *UserStore) createInitialUser() error {
	log.Debug("Creating initial admin user")

	user := UserModel{
		Password:       sql.NullString{String: "admin", Valid: true},
		Name:           sql.NullString{String: "admin", Valid: true},
		Email:          sql.NullString{String: "admin@localhost", Valid: true},
		AuthExternalID: sql.NullString{String: "", Valid: false},
		AuthType:       sql.NullInt32{Int32: int32(AuthStandard), Valid: true},
		Role:           sql.NullInt32{Int32: int32(RoleAdmin), Valid: true},
		Theme:          sql.NullInt32{Int32: int32(ThemeLight), Valid: true},
	}

	_, err := store.Insert(&user)
	return err
}

func (store *UserStore) scanRows(rows *sql.Rows) ([]*User, error) {
	users := []*User{}
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

func (store *UserStore) scanRow(row Scannable) (*User, error) {
	user := new(User)
	timeCreated := int64(0)
	authExternalID := []byte{}
	err := row.Scan(
		&user.ID,
		&timeCreated,
		&user.PasswordHash,
		&user.Name,
		&user.Email,
		&user.AuthType,
		&authExternalID,
		&user.Role,
		&user.Theme,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %s", err)
	}

	user.TimeCreated = time.Unix(timeCreated, 0)
	user.AuthExternalID = string(authExternalID)
	log.Tracef("%+v", user)
	return user, nil
}
