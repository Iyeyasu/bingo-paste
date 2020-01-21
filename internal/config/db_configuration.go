package config

// DBConfiguration contains all database settings
type DBConfiguration struct {
	Driver   string `yaml:"driver"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SSL      string `yaml:"ssl"`
}
