package config

const (
	// AuthStandard users log in using a stored password.
	AuthStandard AuthMode = iota

	// AuthLDAP users log in using their LDAP credentials.
	AuthLDAP = iota
)

const (
	// RoleViewer can view pastes.
	RoleViewer Role = iota

	// RoleEditor can create new pastes.
	RoleEditor = iota

	// RoleAdmin can change site configuration and users.
	RoleAdmin = iota
)

// AuthMode represents the access level of the user.
type AuthMode int

// Role represents the access level of the user.
type Role int

// AuthConfig contains configuration for user authentication.
type AuthConfig struct {
	Enabled        bool     `yaml:"enabled"`
	DefaultMode    AuthMode `yaml:"-"`
	RawDefaultMode string   `yaml:"default_mode"`
	DefaultRole    Role     `yaml:"-"`
	RawDefaultRole string   `yaml:"default_role"`

	Session struct {
		Name         string `yaml:"name"`
		SecureCookie bool   `yaml:"secure_cookie"`
		Store        string `yaml:"store"`
		Redis        struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Password string `yaml:"password"`
			Database int    `yaml:"database"`
		} `yaml:"redis"`
	} `yaml:"session"`

	Standard struct {
		Enabled           bool `yaml:"enabled"`
		AllowRegistration bool `yaml:"allow_registration"`
	} `yaml:"standard"`

	LDAP struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		UseTLS  bool   `yaml:"tls"`
		Base    string `yaml:"base"`
		Bind    struct {
			DN       string `yaml:"dn"`
			Password string `yaml:"password"`
		} `yaml:"bind"`
		Attributes struct {
			UID   string `yaml:"uid"`
			Name  string `yaml:"name"`
			Email string `yaml:"email"`
		} `yaml:"attributes"`
		UserFilter  string `yaml:"user_filter"`
		AdminFilter string `yaml:"admin_filter"`
	} `yaml:"ldap"`
}

// DefaultAuthConfig creates a new AuthConfig with default values.
func DefaultAuthConfig() AuthConfig {
	config := AuthConfig{
		Enabled:        false,
		DefaultMode:    AuthStandard,
		RawDefaultMode: "standard",
		DefaultRole:    RoleEditor,
		RawDefaultRole: "editor",
	}

	config.Session.Name = "session_bingo"
	config.Session.SecureCookie = false
	config.Session.Store = "memory"

	config.Session.Redis.Host = "localhost"
	config.Session.Redis.Port = 6379
	config.Session.Redis.Password = ""
	config.Session.Redis.Database = 0

	config.Standard.Enabled = true
	config.Standard.AllowRegistration = false

	config.LDAP.Enabled = false

	return config
}

func newAuthMode(authMode string) AuthMode {
	switch authMode {
	case "standard":
		return AuthStandard
	case "ldap":
		return AuthLDAP
	default:
		return AuthStandard
	}
}

func newRole(role string) Role {
	switch role {
	case "admin":
		return RoleAdmin
	case "editor":
		return RoleEditor
	case "viewer":
		return RoleViewer
	default:
		return RoleViewer
	}
}

func (mode AuthMode) String() string {
	switch mode {
	case AuthStandard:
		return "Standard"
	case AuthLDAP:
		return "LDAP"
	default:
		return "<invalid_auth_mode>"
	}
}

func (role Role) String() string {
	switch role {
	case RoleAdmin:
		return "Admin"
	case RoleEditor:
		return "Editor"
	case RoleViewer:
		return "Viewer"
	default:
		return "<invalid_theme>"
	}
}
