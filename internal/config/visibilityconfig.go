package config

const (
	// VisibilityUnlisted hides the paste from searches and listings.
	VisibilityUnlisted Visibility = iota

	// VisibilityListed shows the paste in searches and listings.
	VisibilityListed = iota

	// VisibilityPublic enables sharing the paste to unauthenticated users.
	VisibilityPublic = iota
)

// Visibility determins how the paste shows up in searches and listings.
type Visibility int

// VisibilityConfig contains configuration for paste visibility.
type VisibilityConfig struct {
	Enabled    bool       `yaml:"enabled"`
	Default    Visibility `yaml:"-"`
	RawDefault string     `yaml:"default"`
}

// DefaultVisibilityConfig creates a new VisibilityConfig with default values.
func DefaultVisibilityConfig() VisibilityConfig {
	return VisibilityConfig{
		Enabled:    false,
		Default:    VisibilityListed,
		RawDefault: "listed",
	}
}

func newVisibility(visibility string) Visibility {
	switch visibility {
	case "unlisted":
		return VisibilityUnlisted
	case "listed":
		return VisibilityListed
	case "public":
		return VisibilityPublic
	default:
		return VisibilityUnlisted
	}
}
