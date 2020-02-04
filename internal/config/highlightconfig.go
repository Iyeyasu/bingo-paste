package config

import (
	"sort"
	"strings"

	"github.com/alecthomas/chroma/lexers"
)

var mainLanguages = []string{
	"Bash", "C", "C#", "C++",
	"CMake", "CSS", "HTML", "INI",
	"Java", "JavaScript", "JSON", "PHP",
	"PowerShell", "Python", "Python 3",
	"SQL", "TypeScript", "XML", "YAML",
}

// HighlightConfig contains configuration for syntax highlighting.
type HighlightConfig struct {
	Enabled   bool     `yaml:"enabled"`
	Languages []string `yaml:"languages"`
}

// DefaultHighlightConfig creates a new HighlightConfig with default values.
func DefaultHighlightConfig() HighlightConfig {
	return HighlightConfig{
		Enabled:   true,
		Languages: getLanguages(),
	}
}

func isMainLanguage(lang string) bool {
	for _, mainLang := range mainLanguages {
		if lang == mainLang {
			return true
		}
	}
	return false
}

func getLanguages() []string {
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
