package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string    `yaml:"env"`
	Collector Collector `yaml:"collector"`
	Checker   Checker   `yaml:"checker"`
	Scheduler Scheduler `yaml:"scheduler"`
	Database  Database  `yaml:"database"`
}

type Collector struct {
	Timeout time.Duration `yaml:"timeout"`
}

type Checker struct {
	Timeout time.Duration `yaml:"timeout"`
}

type Database struct {
	// yaml > flag > env > panic
	Path    string        `yaml:"path" env:"DATABASE_PATH" required:"true"`
	Timeout time.Duration `yaml:"timeout"`
}

type Scheduler struct {
	StartImmediately bool          `yaml:"start_immediately"`
	Timeout          time.Duration `yaml:"timeout"`
	Interval         time.Duration `yaml:"interval"`
}

func MustLoadConfig() *Config {
	flags := ParseFlags()

	configPath := GetConfigPath(flags)
	config := MustReadConfigFile(configPath)

	databasePath := GetDatabasePath(flags)
	if databasePath != "" {
		config.Database.Path = databasePath
	}

	return config
}

func MustReadConfigFile(configPath string) *Config {
	_, err := os.Stat(configPath)
	if err != nil {
		panic("Config file not found: " + configPath)
	}

	var config Config
	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		panic("Config file is not valid")
	}

	return &config
}

func ParseFlags() map[string]string {
	var databasePath string
	var configPath string

	flag.StringVar(&databasePath, "db", "", "path to database file")
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	return map[string]string{
		"DATABASE_PATH": databasePath,
		"CONFIG_PATH":   configPath,
	}
}

func GetDatabasePath(flags map[string]string) string {
	if flags["DATABASE_PATH"] != "" {
		return flags["DATABASE_PATH"]
	}

	return ""
}

func GetConfigPath(flags map[string]string) string {
	if flags["CONFIG_PATH"] != "" {
		return flags["CONFIG_PATH"]
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	panic("Config path is not set")
}
