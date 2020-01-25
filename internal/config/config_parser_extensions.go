package config

import (
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/lexers"
	log "github.com/sirupsen/logrus"
)

func (parser *configParser) getLogLevel(name string, defaultValue string) log.Level {
	logLevel := parser.getStr(name, defaultValue)
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

func (parser *configParser) getLanguages(name string, defaultValue string) []string {
	var languages []string
	languageSet := parser.getStr(name, defaultValue)

	switch languageSet {
	case "base":
		languages = getBaseLanguages()
	case "all":
		languages = getAllLanguages()
	case "custom":
		log.Debug("Using language set 'custom'")
		languages = parser.getStrArray("extensions.highlight.languages", []string{})
	default:
		log.Fatalf("Failed to parse log file: unknown value 'language_set: %s'", languageSet)
		return nil
	}

	log.Debugf("Using '%s' language set (%d languages)", languageSet, len(languages))
	log.Debugf("Used languages are %v", languages)
	return languages
}

func (parser *configParser) getDurations(name string, defaultValue []int) []time.Duration {
	dur := parser.getIntArray(name, defaultValue)
	durations := make([]time.Duration, len(dur), len(dur))
	for i := 0; i < len(dur); i++ {
		durations[i] = time.Duration(dur[i]) * time.Minute
	}

	log.Debugf("Using %d expiry durations", len(durations))
	log.Debugf("Used expiry durations are %v", durations)
	return durations
}

func getBaseLanguages() []string {
	return []string{
		"Bash", "C", "C#", "C++", "Clojure",
		"CMake", "CSS", "Dart", "Go", "GLSL",
		"Elixir", "Haskell", "HTML", "HTTP", "INI",
		"Java", "JavaScript", "JSON", "Kotlin", "Objective-C",
		"Perl", "PHP", "PowerShell", "Python", "Python 3",
		"R", "Ruby", "Rust", "Scala", "SQL",
		"Swift", "TypeScript", "VBA", "XML", "YAML",
	}
}

func getAllLanguages() []string {
	langCount := lexers.Registry.Lexers.Len()
	languages := make([]string, 0, langCount)
	for i := 0; i < langCount; i++ {
		name := lexers.Registry.Lexers[i].Config().Name
		if name != "plaintext" {
			languages = append(languages, name)
		}
	}

	sort.Slice(languages, func(i, j int) bool {
		return strings.ToLower(languages[i]) < strings.ToLower(languages[j])
	})

	return languages
}
