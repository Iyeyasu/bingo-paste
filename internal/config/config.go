package config

import (
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

var (
	conf *Config
)

// Config contains all settings.
type Config struct {
	LogLevel log.Level

	API struct {
		URI string
	}

	Appearance struct {
		Title string
		Theme string
		Icon  string
	}

	Cache struct {
		Enabled  bool
		Host     string
		Port     int
		Password string
		Database int
	}

	Database struct {
		Driver   string
		Username string
		Password string
		Database string
		Host     string
		Port     int
		SSL      string
	}

	Extensions struct {
		Visibility struct {
			Enabled bool
		}

		Expiry struct {
			Enabled   bool
			Durations []time.Duration
		}

		Highlight struct {
			Enabled   bool
			Languages []string
		}
	}
}

// Get returns a default configuration.
func Get() *Config {
	if conf == nil {
		conf = NewConfig("/bingo/bingo.yml")
	}
	return conf
}

// NewConfig returns a new configuration read from the given file.
func NewConfig(filename string) *Config {
	loadLogLevel(filename)

	parser, err := newConfigParser(filename)
	if err != nil {
		log.Fatalf(err.Error())
	}

	conf := new(Config)
	conf.LogLevel = parser.getLogLevel("log_level", "info")
	conf.API.URI = parser.getStr("api.uri", "/api/v1")

	conf.Appearance.Title = parser.getStr("appearance.title", "Bingo")
	conf.Appearance.Theme = parser.getStr("appearance.theme", "dark")
	conf.Appearance.Icon = parser.getStr("appearance.icon", "")

	conf.Cache.Enabled = parser.hasOption("redis")
	conf.Cache.Host = parser.getStr("redis.host", "localhost")
	conf.Cache.Port = parser.getInt("redis.port", 6379)
	conf.Cache.Password = parser.getStr("redis.password", "")
	conf.Cache.Database = parser.getInt("redis.database", 0)

	conf.Database.Driver = parser.getStr("db.driver", "sqlite3")
	conf.Database.Username = parser.getStr("db.username", "")
	conf.Database.Password = parser.getStr("db.password", "")
	conf.Database.Database = parser.getStr("db.database", "")
	conf.Database.Host = parser.getStr("db.host", "localhost")
	conf.Database.SSL = parser.getStr("db.ssl", "required")

	switch conf.Database.Driver {
	case "postgres":
		conf.Database.Port = parser.getInt("db.port", 5432)
	case "mysql":
		conf.Database.Port = parser.getInt("db.port", 3306)
	}

	conf.Extensions.Visibility.Enabled = parser.getBool("extensions.visibility.enabled", false)
	if conf.Extensions.Visibility.Enabled {
		log.Debug("Visibility extension enabled")
	}

	conf.Extensions.Highlight.Enabled = parser.getBool("extensions.highlight.enabled", false)
	if conf.Extensions.Highlight.Enabled {
		log.Debug("Syntax highlight extension enabled")
		conf.Extensions.Highlight.Languages = parser.getLanguages()
	}

	conf.Extensions.Expiry.Enabled = parser.getBool("extensions.expiry.enabled", false)
	if conf.Extensions.Expiry.Enabled {
		log.Debug("Expiry extension enabled")
		defaultDurations := []int{10, 60, 1440, 10080, 43200, 525600}
		conf.Extensions.Expiry.Durations = parser.getDurations("extensions.expiry.durations", defaultDurations)
	}

	log.Tracef("Parsed configuration: %+v", conf)

	return conf
}

// Load log level before everything else so that we get logging for config
// file reading as well.
func loadLogLevel(filename string) *Config {
	parser, err := newConfigParser(filename)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.SetLevel(parser.getLogLevel("log_level", "info"))
	return conf
}
