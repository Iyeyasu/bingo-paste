package config

// DatabaseConfig contains configuration for database connection.
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SSL      string `yaml:"ssl"`
}

// DefaultDatabaseConfig creates a new DatabaseConfig with default values.
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:   "postgres",
		Username: "",
		Password: "",
		Database: "",
		Host:     "localhost",
		Port:     5432,
		SSL:      "required",
	}
}
