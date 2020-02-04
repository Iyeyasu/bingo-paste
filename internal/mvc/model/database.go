package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bingo/internal/config"
	"bingo/internal/util/log"

	// Postgresql driver
	_ "github.com/lib/pq"
)

var (
	dbConnectionRetries     = 45
	dbConnectionTimeout     = 60
	dbMaxOpenConnections    = 20
	dbMaxIdleConnections    = 10
	dbMaxConnectionLifetime = 5 * time.Minute
)

// NewDatabase returns a new SQL database connection.
func NewDatabase() *sql.DB {
	driver, connStr, err := getDataSource()
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	log.Infof("Opening %s database", driver)
	log.Debugf("Database connection string: %s", connStr)
	db, err := sql.Open(driver, connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	err = pollDatabase(db)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	configureDatabase(db)
	createPseudoEncrypt(db)
	return db
}

func pollDatabase(db *sql.DB) error {
	log.Infof("Trying to connect to database for %d seconds", dbConnectionRetries)

	for i := 0; i <= dbConnectionRetries; i++ {
		err := db.Ping()
		if err == nil {
			log.Info("Connected to database")
			return nil
		}

		log.Info(err)
		time.Sleep(time.Second)
	}
	return errors.New("failed to connect to database")
}

func configureDatabase(db *sql.DB) {
	log.Debug("Configuring database")

	log.Debugf("Setting max open connections to %d", dbMaxOpenConnections)
	db.SetMaxOpenConns(dbMaxOpenConnections)
	log.Debugf("Setting max idle connections to %d", dbMaxIdleConnections)
	db.SetMaxIdleConns(dbMaxIdleConnections)
	log.Debugf("Setting connection max lifetime to %d", int64(dbMaxConnectionLifetime.Seconds()))
	db.SetConnMaxLifetime(dbMaxConnectionLifetime)
}

func getDataSource() (string, string, error) {
	driver := config.Get().Database.Driver
	connStr := ""

	switch driver {
	case "postgres":
		connStr = getPostgresConnectionString()
	default:
		return "", "", fmt.Errorf("invalid database driver '%s'", driver)
	}

	return driver, connStr, nil
}

func getPostgresConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=%d",
		config.Get().Database.Username,
		config.Get().Database.Password,
		config.Get().Database.Host,
		config.Get().Database.Port,
		config.Get().Database.Database,
		config.Get().Database.SSL,
		dbConnectionTimeout)
}

// https://stackoverflow.com/questions/12761346/pseudo-encrypt-function-in-plpgsql-that-takes-bigint/12761795#12761795
// Creates a function that maps big integers to another seemingly random big integer.
// Used to make sure that object ids for are seemingly random.
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
