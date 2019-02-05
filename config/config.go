package config

import (
	"github.com/BurntSushi/toml"
)

// Config all settings
type Config struct {
	Server   Server `toml:"server"`
	DBMaster DB     `toml:"dbm"`
	DBSlave  DB     `toml:"dbs"`
}

// Server port
type Server struct {
	RedirectURL string `toml:"redirect_url"`
	Port        string `toml:"port"`
}

// DB database structure
type DB struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

// New Config
func New(config *Config, configPath string) error {
	_, err := toml.DecodeFile(configPath, config)
	return err
}
