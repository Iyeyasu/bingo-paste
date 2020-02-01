package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/util"
	log "github.com/sirupsen/logrus"
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
	err := store.query.countStmt.QueryRow().Scan(&count)
	if err != nil {
		log.Fatalf("Failed to get user count: %s", err)
	}
	return count
}

// Insert inserts a new user to the database.
func (store *UserStore) Insert(userModel *UserModel) (*User, error) {
	log.Debug("Inserting new user to database")

	passwordHash, err := util.HashPassword([]byte(userModel.Password.String))
	if err != nil {
		return nil, err
	}

	timeCreated := time.Now().Unix()
	row := store.query.insertStmt.QueryRow(
		timeCreated,
		passwordHash,
		userModel.Name,
		userModel.Mail,
		userModel.AuthType,
		userModel.AuthExternalID,
		userModel.Role,
		userModel.Theme,
	)

	return store.scanRow(row)
}

// Select returns the user with the given id from the database.
func (store *UserStore) Select(id int64) (*User, error) {
	log.Debugf("Retrieving user %d from database", id)

	row := store.query.selectStmt.QueryRow(id)
	return store.scanRow(row)
}

// SelectMultiple returns a slice of public users sorted by their creation time.
func (store *UserStore) SelectMultiple(limit int64, offset int64) ([]*User, error) {
	log.Debugf("Retrieving %d public users starting from user number %d from database", limit, offset)

	rows, err := store.query.selectMultipleStmt.Query(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %s", err)
	}

	return store.scanRows(rows)
}

// Update Updates an existing user in the database.
func (store *UserStore) Update(userModel *UserModel) (*User, error) {
	log.Debug("Updating existing user in the database")

	row := store.query.updateStmt.QueryRow(
		userModel.ID,
		nil,
		userModel.Name,
		userModel.Mail,
		userModel.AuthType,
		userModel.AuthExternalID,
		userModel.Role,
		userModel.Theme,
	)

	return store.scanRow(row)
}

// Delete deletes the user with the given id from the database.
func (store *UserStore) Delete(id int64) error {
	log.Debugf("Deleting user %d from database", id)

	_, err := store.query.deleteStmt.Exec(id)
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
		Mail:           sql.NullString{String: "admin@localhost", Valid: true},
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
		&user.Mail,
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
	return user, nil
}
