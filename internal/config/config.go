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
	Path    string        `yaml:"path"`
	Timeout time.Duration `yaml:"timeout"`
}

type Scheduler struct {
	Interval time.Duration `yaml:"interval"`
}

func MustLoad(configPath string) *Config {
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

func GetConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()
	if configPath != "" {
		return configPath
	}

	panic("Config path is not set")
}
