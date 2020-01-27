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
	dbConnectionRetries   = 45
	connectionTimeout     = 60
	maxOpenConnections    = 20
	maxIdleConnections    = 10
	maxConnectionLifetime = 5 * time.Minute
)

// Scannable is a common interface for sql.Row and sql.Rows.
type Scannable interface {
	Scan(dest ...interface{}) error
}

// NewDatabase opens a new SQL connection.
func NewDatabase() *sql.DB {
	driver, connStr, err := getDataSource()

	log.Infof("Opening %s database", driver)
	log.Debugf("Database connection string: %s", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

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
	log.Infof("Trying to connect to database for %d seconds", connectionRetries)

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

	log.Debugf("Setting max open connections to %d", maxOpenConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	log.Debugf("Setting max idle connections to %d", maxIdleConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	log.Debugf("Setting connection max lifetime to %d", int64(maxConnectionLifetime.Seconds()))
	db.SetConnMaxLifetime(maxConnectionLifetime)
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
		connectionTimeout)
}
