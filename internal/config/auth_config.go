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
	Registration   bool     `yaml:"registration"`
	SessionName    string   `yaml:"session_name"`
	SecureCookie   bool     `yaml:"secure_cookie"`
}

// DefaultAuthConfig creates a new AuthConfig with default values.
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		Enabled:        false,
		DefaultMode:    AuthStandard,
		RawDefaultMode: "standard",
		DefaultRole:    RoleEditor,
		RawDefaultRole: "editor",
		Registration:   false,
		SessionName:    "session_bingo",
		SecureCookie:   false,
	}
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
		return "Editor"
	default:
		return "<invalid_theme>"
	}
}
