package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/config"

	// Postgres driver
	_ "github.com/lib/pq"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"

	// Sqllite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// OpenDB opens the connection to the database.
func OpenDB() (*sql.DB, error) {
	retries := 45
	driver, connStr := getDataSource()
	log.Printf("Connecting to %s database. Polling for %d seconds...", driver, retries)

	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	for i := 0; i <= retries; i++ {
		err = db.Ping()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		log.Println("Connected to database.")
		return db, nil
	}

	return nil, errors.New("couldn't open db connection")
}

func getDataSource() (string, string) {
	conf := config.ReadFile("/bingo/bingo.yml")
	driver := conf.Database.Driver
	timeout := 60

	switch driver {
	case "sqlite3":
		return driver, "file:bingo.db?cache=shared&mode=memory"
	case "postgres":
		return driver, fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s&connect_timeout=%d",
			conf.Database.Username,
			conf.Database.Password,
			conf.Database.Host,
			conf.Database.Port,
			conf.Database.Database,
			conf.Database.SSL,
			timeout)
	case "mysql":
		return driver, fmt.Sprintf(
			"%s:%s@%s(%s:%d)/%s?tls=%s&timeout=%d",
			conf.Database.Username,
			conf.Database.Password,
			"tcp",
			conf.Database.Host,
			conf.Database.Port,
			conf.Database.Database,
			conf.Database.SSL,
			timeout)
	default:
		log.Fatalf("Invalid database driver '%s'\n", driver)
		return "", ""
	}
}
