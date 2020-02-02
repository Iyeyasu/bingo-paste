package model

import (
	"database/sql"

	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

type userQuery struct {
	count       *sql.Stmt
	insert      *sql.Stmt
	findByID    *sql.Stmt
	findByName  *sql.Stmt
	findByEmail *sql.Stmt
	findRange   *sql.Stmt
	update      *sql.Stmt
	delete      *sql.Stmt
}

func newUserQuery(db *sql.DB) *userQuery {
	query := new(userQuery)
	createTable(db)
	query.count = createCountStatement(db)
	query.findByID = createFindByIDStatement(db)
	query.findByName = createFindByNameStatement(db)
	query.findByEmail = createFindByEmailStatement(db)
	query.findRange = createFindRangeStatement(db)
	query.insert = createInsertStatement(db)
	query.update = createUpdateStatement(db)
	query.delete = createDeleteStatement(db)
	return query
}

func createTable(db *sql.DB) {
	log.Debug("Creating table 'users'")

	q := `
		CREATE SEQUENCE IF NOT EXISTS users_id_seq AS bigint;

		CREATE TABLE IF NOT EXISTS users (
			id 					bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('users_id_seq')),
			time_created_sec 	bigint NOT NULL,
			password_hash 		text,
			name 				text NOT NULL,
			email 				text NOT NULL,
			auth_type 			int NOT NULL,
			auth_external_id 	text,
			role 				int NOT NULL,
			theme 				int NOT NULL
		);

		ALTER SEQUENCE users_id_seq OWNED BY users.id
	`

	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table 'users': %s", err)
	}
}

func createCountStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT COUNT(*) FROM USERS")
	if err != nil {
		log.Fatalf("Failed to create count users statement: %s", err)
	}
	return stmt
}

func createFindByIDStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to create find user by id statement: %s", err)
	}
	return stmt
}

func createFindByNameStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT * FROM users WHERE name ILIKE $1")
	if err != nil {
		log.Fatalf("Failed to create find user by name user statement: %s", err)
	}
	return stmt
}

func createFindByEmailStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT * FROM users WHERE email ILIKE $1")
	if err != nil {
		log.Fatalf("Failed to create find user by email user statement: %s", err)
	}
	return stmt
}

func createFindRangeStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("SELECT * FROM users ORDER BY name DESC, id ASC LIMIT $1 OFFSET $2")
	if err != nil {
		log.Fatalf("Failed to create select user list statement: %s", err)
	}
	return stmt
}

func createInsertStatement(db *sql.DB) *sql.Stmt {
	q := `
		INSERT INTO users (
				time_created_sec,
				password_hash,
				name,
				email,
				auth_type,
				auth_external_id,
				role,
				theme)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
	`

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to create insert user statement: %s", err)
	}
	return stmt
}

func createUpdateStatement(db *sql.DB) *sql.Stmt {
	q := `
		UPDATE users
		SET
			password_hash 		= COALESCE($2, password_hash),
			name 				= COALESCE($3, name),
			email 				= COALESCE($4, email),
			auth_type 			= COALESCE($5, auth_type),
			auth_external_id 	= COALESCE($6, auth_external_id),
			role 				= COALESCE($7, role),
			theme 				= COALESCE($8, theme)
		WHERE id = $1
		RETURNING *
	`

	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalf("Failed to create update user statement: %s", err)
	}
	return stmt
}

func createDeleteStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to create delete user statement: %s", err)
	}
	return stmt
}
