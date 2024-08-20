package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env"`
	Collector  Collector  `yaml:"collector"`
	Checker    Checker    `yaml:"checker"`
	Scheduler  Scheduler  `yaml:"scheduler"`
	Database   Database   `yaml:"database"`
	Mail       Mail       `yaml:"mail"`
	Rabbit     Rabbit     `yaml:"rabbit"`
	GoogleAuth GoogleAuth `yaml:"google_auth"`
	JWT        JWT        `yaml:"jwt"`
}

type Mail struct {
	From    string        `yaml:"from" required:"true"`
	Mail    string        `yaml:"mail" required:"true"`
	Pass    string        `yaml:"pass" required:"true"`
	Secure  bool          `yaml:"secure" default:"false"`
	Addr    string        `yaml:"addr" required:"true"`
	Port    int           `yaml:"port" default:"587"`
	Timeout time.Duration `yaml:"timeout"`
}

type GoogleAuth struct {
	ClientId     string         `yaml:"client_id" required:"true"`
	ClientSecret string         `yaml:"client_secret" required:"true"`
	RedirectUrl  string         `yaml:"redirect_url" required:"true"`
	Scope        []string       `yaml:"scope" required:"true"`
	RedirectTo   string         `yaml:"redirect_to" required:"true"`
	Timeout      time.Duration  `yaml:"timeout"`
}

type JWT struct {
	Secret          string        `yaml:"secret" required:"true"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" default:"24h"`
	Timeout         time.Duration `yaml:"timeout"`
}

type Rabbit struct {
	Host    string        `yaml:"host" required:"true"`
	Port    int           `yaml:"port" default:"5672"`
	User    string        `yaml:"user" required:"true"`
	Pass    string        `yaml:"pass" required:"true"`
	Timeout time.Duration `yaml:"timeout"`
}

type Collector struct {
	Timeout time.Duration `yaml:"timeout"`
}

type Checker struct {
	Timeout time.Duration `yaml:"timeout"`
}

type Database struct {
	//  flag > yaml > panic
	Path    string        `yaml:"path" required:"true"`
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

	if flags["DATABASE_PATH"] != "" {
		config.Database.Path = flags["DATABASE_PATH"]
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
