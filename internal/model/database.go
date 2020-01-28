package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	log "github.com/sirupsen/logrus"

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

// Scannable is a common interface for sql.Row and sql.Rows.
type Scannable interface {
	Scan(dest ...interface{}) error
}

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
