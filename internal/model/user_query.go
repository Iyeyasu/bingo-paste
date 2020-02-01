package model

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	userColumns = `
		time_created_sec,
		password_hash,
		name,
		mail,
		auth_type,
		auth_external_id,
		role,
		theme
		`
)

type userQuery struct {
	countStmt          *sql.Stmt
	insertStmt         *sql.Stmt
	selectStmt         *sql.Stmt
	selectMultipleStmt *sql.Stmt
	updateStmt         *sql.Stmt
	deleteStmt         *sql.Stmt
}

func newUserQuery(db *sql.DB) *userQuery {
	query := new(userQuery)
	query.createTable(db)
	query.countStmt = query.createCountStatement(db)
	query.insertStmt = query.createInsertStatement(db)
	query.selectStmt = query.createSelectStatement(db)
	query.selectMultipleStmt = query.createSelectMultipleStatement(db)
	query.updateStmt = query.createUpdateStatement(db)
	query.deleteStmt = query.createDeleteStatement(db)
	return query
}

func (store *userQuery) createTable(db *sql.DB) {
	log.Debug("Creating table 'users'")

	q := "CREATE SEQUENCE IF NOT EXISTS users_id_seq AS bigint"
	_, err := db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create sequence 'users_id_seq': %s", err)
	}

	q = `
	CREATE TABLE IF NOT EXISTS users (
		id 					bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('users_id_seq')),
		time_created_sec 	bigint NOT NULL,
		password_hash 		text,
		name 				text NOT NULL,
		mail 				text NOT NULL,
		auth_type 			int NOT NULL,
		auth_external_id 	text,
		role 				int NOT NULL,
		theme 				int NOT NULL
	)
	`
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to create table 'users': %s", err)
	}

	q = "ALTER SEQUENCE users_id_seq OWNED BY users.id"
	_, err = db.Exec(q)
	if err != nil {
		log.Fatalf("Failed to assign sequence 'users_id_seq': %s", err)
	}
}

func (store *userQuery) createCountStatement(db *sql.DB) *sql.Stmt {
	query := "SELECT COUNT(*) FROM USERS"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to create count users statement: %s", err)
	}
	return stmt
}

func (store *userQuery) createInsertStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
	INSERT INTO users (%s)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, %s
	`, userColumns, userColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to create insert user statement: %s", err)
	}
	return stmt
}

func (store *userQuery) createSelectStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
	SELECT id, %s
	FROM users
	WHERE id = $1
	`, userColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to create select user statement: %s", err)
	}
	return stmt
}

func (store *userQuery) createSelectMultipleStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
	SELECT id, %s
	FROM users
	ORDER BY name DESC, id ASC
	LIMIT $1 OFFSET $2
	`, userColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to create select user list statement: %s", err)
	}
	return stmt
}

func (store *userQuery) createUpdateStatement(db *sql.DB) *sql.Stmt {
	query := fmt.Sprintf(`
	UPDATE users
	SET
		password_hash 		= COALESCE($2, password_hash),
		name 				= COALESCE($3, name),
		mail 				= COALESCE($4, mail),
		auth_type 			= COALESCE($5, auth_type),
		auth_external_id 	= COALESCE($6, auth_external_id),
		role 				= COALESCE($7, role),
		theme 				= COALESCE($8, theme)
	WHERE id = $1
	RETURNING id, %s
	`, userColumns)

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to create update user statement: %s", err)
	}
	return stmt
}

func (store *userQuery) createDeleteStatement(db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		log.Fatalf("Failed to create delete user statement: %s", err)
	}
	return stmt
}
