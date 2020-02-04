package store

import (
	"database/sql"
	"time"

	"bingo/internal/mvc/model"
	"bingo/internal/util/auth"
	"bingo/internal/util/log"

	"github.com/jmoiron/sqlx"
)

// UserStore is the store for users.
type UserStore struct {
	Database *sqlx.DB
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *sqlx.DB) *UserStore {
	log.Debug("Initializing user store")
	store := new(UserStore)
	store.Database = db
	return store
}

// Count returns the number of users.
func (store *UserStore) Count() int64 {
	log.Debugf("Counting number of users")

	var count int64
	store.Database.Get(&count, "SELECT COUNT(*) FROM users")
	return count
}

// FindByID returns the user with the given id from the database.
func (store *UserStore) FindByID(id int64) (*model.User, error) {
	log.Debugf("Retrieving user %d from database", id)

	query := `
		SELECT id, time_created, name, email, password_hash, auth_mode, auth_external_id, role, theme
		FROM users
		WHERE id = $1
		`

	user := new(model.User)
	err := store.Database.Get(user, query, id)
	return user, err
}

// FindByName returns the user with the given name from the database.
func (store *UserStore) FindByName(name string) (*model.User, error) {
	log.Debugf("Retrieving user with name '%s' from database", name)

	query := `
		SELECT id, time_created, name, email, password_hash, auth_mode, auth_external_id, role, theme
		FROM users
		WHERE lower(name) = lower($1)
		`

	user := new(model.User)
	err := store.Database.Get(user, query, name)
	return user, err
}

// FindByEmail returns the user with the given email from the database.
func (store *UserStore) FindByEmail(email string) (*model.User, error) {
	log.Debugf("Retrieving user with mail '%s' from database", email)

	query := `
		SELECT id, time_created, name, email, password_hash, auth_mode, auth_external_id, role, theme
		FROM users
		WHERE email ILIKE $1
		`

	user := new(model.User)
	err := store.Database.Get(user, query, email)
	return user, err
}

// FindRange returns a slice of public users sorted by their creation time.
func (store *UserStore) FindRange(limit int64, offset int64) ([]model.User, error) {
	log.Debugf("Retrieving %d public users starting from user number %d from database", limit, offset)

	query := `
		SELECT id, time_created, name, email, password_hash, auth_mode, auth_external_id, role, theme
		FROM users
		ORDER BY role DESC, name ASC, id ASC
		LIMIT $1 OFFSET $2
		`

	users := []model.User{}
	err := store.Database.Select(&users, query, limit, offset)
	return users, err
}

// Delete deletes the user with the given id from the database.
func (store *UserStore) Delete(id int64) error {
	log.Debugf("Deleting user %d from database", id)

	_, err := store.Database.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// Insert inserts a new user to the database.
func (store *UserStore) Insert(userTmpl *model.UserTemplate) (*model.User, error) {
	log.Debug("Inserting new user to database")
	log.Tracef("%+v", userTmpl)

	query := `
		INSERT INTO users (
				time_created,
				password_hash,
				name,
				email,
				auth_mode,
				auth_external_id,
				role,
				theme)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
	`

	passwordHash, err := auth.HashPassword(userTmpl.Password.String)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	err = store.Database.QueryRowx(
		query,
		time.Now().UTC(),
		passwordHash,
		userTmpl.Name,
		userTmpl.Email,
		userTmpl.AuthMode,
		userTmpl.AuthExternalID,
		userTmpl.Role,
		userTmpl.Theme,
	).StructScan(user)

	return user, err
}

// Update Updates an existing user in the database.
func (store *UserStore) Update(userTmpl *model.UserTemplate) (*model.User, error) {
	log.Debug("Updating existing user in the database")
	log.Tracef("%+v", userTmpl)

	query := `
		UPDATE users
		SET
			password_hash 		= COALESCE($2, password_hash),
			name 				= COALESCE($3, name),
			email 				= COALESCE($4, email),
			auth_mode 			= COALESCE($5, auth_mode),
			auth_external_id 	= COALESCE($6, auth_external_id),
			role 				= COALESCE($7, role),
			theme 				= COALESCE($8, theme)
		WHERE id = $1
		RETURNING *
	`

	var passwordHash sql.NullString
	if userTmpl.Password.Valid {
		hash, err := auth.HashPassword(userTmpl.Password.String)
		if err != nil {
			return nil, err
		}

		passwordHash.String = hash
		passwordHash.Valid = true
	}

	user := new(model.User)
	err := store.Database.QueryRowx(
		query,
		userTmpl.ID,
		passwordHash,
		userTmpl.Name,
		userTmpl.Email,
		userTmpl.AuthMode,
		userTmpl.AuthExternalID,
		userTmpl.Role,
		userTmpl.Theme,
	).StructScan(user)

	return user, err
}

func (store *UserStore) createTable() {
	query := `
		CREATE SEQUENCE IF NOT EXISTS users_id_seq AS bigint;

		CREATE TABLE IF NOT EXISTS users (
			id 					bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('users_id_seq')),
			time_created 	    timestamptz NOT NULL,
			password_hash 		char(60),
			name 				varchar(32) NOT NULL,
			email 				varchar(254) NOT NULL,
			auth_mode 			int NOT NULL,
			auth_external_id 	varchar(32),
			role 				int NOT NULL,
			theme 				int NOT NULL
		);

		ALTER SEQUENCE users_id_seq OWNED BY users.id;
		CREATE UNIQUE INDEX IF NOT EXISTS name_lower_idx ON users (lower(name));
	`

	_, err := store.Database.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table 'users': %s", err)
	}
}
