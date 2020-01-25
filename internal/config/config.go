package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	conf *Config
)

// Config contains all settings.
type Config struct {
	LogLevel log.Level

	Appearance struct {
		Title string
		Theme string
		Icon  string
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

	Cache struct {
		Host     string
		Port     string
		Password string
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

// Get returns the configuration.
func Get() *Config {
	if conf == nil {
		conf = newConfig("/bingo/bingo.yml")
	}
	return conf
}

func newConfig(filename string) *Config {
	parser := newConfigParser(filename)
	conf := new(Config)

	conf.LogLevel = parser.getLogLevel("log_level", "info")

	conf.Appearance.Title = parser.getStr("appearance.title", "Bingo")
	conf.Appearance.Theme = parser.getStr("appearance.theme", "dark")
	conf.Appearance.Icon = parser.getStr("appearance.icon", "")

	conf.Database.Driver = parser.getStr("db.driver", "sqlite3")
	conf.Database.Username = parser.getStr("db.username", "")
	conf.Database.Password = parser.getStr("db.password", "")
	conf.Database.Database = parser.getStr("db.database", "")
	conf.Database.Host = parser.getStr("db.host", "localhost")
	conf.Database.Port = parser.getInt("db.port", 0)
	conf.Database.SSL = parser.getStr("db.ssl", "required")

	conf.Extensions.Visibility.Enabled = parser.getBool("extensions.visibility.enabled", false)
	if conf.Extensions.Visibility.Enabled {
		log.Debug("Visibility extension enabled")
	}

	conf.Extensions.Highlight.Enabled = parser.getBool("extensions.highlight.enabled", false)
	if conf.Extensions.Highlight.Enabled {
		log.Debug("Syntax highlight extension enabled")
		conf.Extensions.Highlight.Languages = parser.getLanguages("extensions.highlight.languageSet", "base")
	}

	conf.Extensions.Expiry.Enabled = parser.getBool("extensions.expiry.enabled", false)
	if conf.Extensions.Expiry.Enabled {
		log.Debug("Expiry extension enabled")
		defaultDurations := []int{0, 10, 60, 1440, 10800, 43200, 525600}
		conf.Extensions.Expiry.Durations = parser.getDurations("extensions.expiry.durations", defaultDurations)
	}

	log.Tracef("Parsed configuration: %+v", conf)

	return conf
}
