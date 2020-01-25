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

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// Sqllite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	connectionRetries     = 45
	connectionTimeout     = 60
	maxOpenConnections    = 20
	maxIdleConnections    = 10
	maxConnectionLifetime = 5 * time.Minute
)

// Database contains all different stores and the SQL connection instance.
type Database struct {
	Pastes *PasteStore
}

// NewDatabase opens a new database connection and initializes stores.
func NewDatabase() *Database {
	driver, connStr, err := getDataSource()
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	db, err := openDatabase(driver, connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	err = pollDatabase(db)
	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	configureDatabase(db)

	database := new(Database)
	database.createStores(db)
	return database
}

func (database *Database) createStores(db *sql.DB) {
	database.Pastes = NewPasteStore(db)
}

func openDatabase(driver string, connStr string) (*sql.DB, error) {
	log.Infof("Opening %s database", driver)
	log.Debugf("Database connection string: %s", connStr)

	return sql.Open(driver, connStr)
}

func pollDatabase(db *sql.DB) error {
	log.Infof("Trying to connect to database for %d seconds", connectionRetries)

	for i := 0; i <= connectionRetries; i++ {
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
	case "sqlite3":
		connStr = "file:sqlite.db?cache=shared&mode=memory"
	case "postgres":
		connStr = getPostgresConnectionString()
	case "mysql":
		connStr = getMySQLConnectionString()
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

func getMySQLConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@%s(%s:%d)/%s?tls=%s&timeout=%d",
		config.Get().Database.Username,
		config.Get().Database.Password,
		"tcp",
		config.Get().Database.Host,
		config.Get().Database.Port,
		config.Get().Database.Database,
		config.Get().Database.SSL,
		connectionTimeout)
}
