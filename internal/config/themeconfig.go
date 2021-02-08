package config

import "log"

const (
	// ThemeLight represents a dark GUI theme.
	ThemeLight Theme = iota

	// ThemeDark represents a light GUI theme.
	ThemeDark = iota
)

// Theme is the theme used by the user.
type Theme int

// ThemeConfig contains configuration for the web GUI appearance.
type ThemeConfig struct {
	Title      string `yaml:"title"`
	Default    Theme  `yaml:"-"`
	RawDefault string `yaml:"default"`
	Icon       string `yaml:"icon"`
}

// DefaultThemeConfig creates a new ThemeConfig with default values.
func DefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		Title:      "Pastebin",
		Default:    ThemeLight,
		RawDefault: "light",
		Icon:       "",
	}
}

func newTheme(theme string) Theme {
	switch theme {
	case "light":
		return ThemeLight
	case "dark":
		return ThemeDark
	default:
		log.Fatalf("Failed to parse log file: unknown value 'theme.default: %s'", theme)
		return ThemeLight
	}
}

func (theme Theme) String() string {
	switch theme {
	case ThemeLight:
		return "Light"
	case ThemeDark:
		return "Dark"
	default:
		return "<invalid_theme>"
	}
}
