package config

import (
	"io/ioutil"
	"time"

	"bingo/internal/util/log"

	"gopkg.in/yaml.v2"
)

var conf *Config

// Config contains all settings.
type Config struct {
	Host           string           `yaml:"host"`
	Port           int              `yaml:"port"`
	LogLevel       log.Level        `yaml:"-"`
	RawLogLevel    string           `yaml:"log_level"`
	Authentication AuthConfig       `yaml:"auth"`
	Database       DatabaseConfig   `yaml:"db"`
	Expiry         ExpiryConfig     `yaml:"expiry"`
	Highlight      HighlightConfig  `yaml:"highlight"`
	Theme          ThemeConfig      `yaml:"theme"`
	Visibility     VisibilityConfig `yaml:"visibility"`
}

// Get returns a default configuration.
func Get() *Config {
	return conf
}

// NewConfig returns a new configuration read from the given file.
func Load(filename string) *Config {
	log.Infof("Reading configuration file '%s'", filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read config file '%s': %s", filename, err)
	}

	conf = NewDefaultConfig()
	err = yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		log.Fatalf("failed to parse config file '%s': %s", filename, err)
	}

	conf.LogLevel = newLogLevel(conf.RawLogLevel)
	conf.Theme.Default = newTheme(conf.Theme.RawDefault)
	conf.Visibility.Default = newVisibility(conf.Visibility.RawDefault)
	conf.Authentication.DefaultMode = newAuthMode(conf.Authentication.RawDefaultMode)
	conf.Authentication.DefaultRole = newRole(conf.Authentication.RawDefaultRole)

	if !conf.Expiry.Enabled {
		conf.Expiry.Durations = []time.Duration{}
	} else {
		conf.Expiry.Durations = newDurations(conf.Expiry.RawDurations)
	}

	if !conf.Highlight.Enabled {
		conf.Highlight.Languages = []string{}
	}

	log.SetLevel(conf.LogLevel)
	log.Tracef("Parsed configuration: %+v", conf)

	return conf
}

// NewDefaultConfig creates a new configuration with default values.
func NewDefaultConfig() *Config {
	conf = new(Config)
	conf.Host = "0.0.0.0"
	conf.Port = 80
	conf.RawLogLevel = "info"
	conf.Authentication = DefaultAuthConfig()
	conf.Database = DefaultDatabaseConfig()
	conf.Expiry = DefaultExpiryConfig()
	conf.Highlight = DefaultHighlightConfig()
	conf.Theme = DefaultThemeConfig()
	conf.Visibility = DefaultVisibilityConfig()
	return conf
}

func newLogLevel(logLevel string) log.Level {
	switch logLevel {
	case "panic":
		return log.PanicLevel
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		log.Fatalf("Failed to parse log file: unknown value 'log_level: %s'", logLevel)
		return log.TraceLevel
	}
}
