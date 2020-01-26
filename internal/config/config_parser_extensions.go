package config

import (
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/lexers"
	log "github.com/sirupsen/logrus"
)

var (
	mainLanguages = []string{
		"Bash", "C", "C#", "C++",
		"CMake", "CSS", "HTML", "INI",
		"Java", "JavaScript", "JSON", "PHP",
		"PowerShell", "Python", "Python 3",
		"SQL", "TypeScript", "XML", "YAML",
	}
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

func (parser *configParser) getLanguages() []string {
	langCount := lexers.Registry.Lexers.Len()
	languages := make([]string, 0, langCount)
	for i := 0; i < langCount; i++ {
		name := lexers.Registry.Lexers[i].Config().Name
		if name != "plaintext" {
			languages = append(languages, name)
		}
	}

	sort.Slice(languages, func(i, j int) bool {
		iIsMain := isMainLanguage(languages[i])
		jIsMain := isMainLanguage(languages[j])
		if iIsMain && !jIsMain {
			return true
		} else if !iIsMain && jIsMain {
			return false
		}
		return strings.ToLower(languages[i]) < strings.ToLower(languages[j])
	})

	return languages
}

func isMainLanguage(lang string) bool {
	for _, mainLang := range mainLanguages {
		if lang == mainLang {
			return true
		}
	}
	return false
}
