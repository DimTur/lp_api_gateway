package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer HTTPServer    `yaml:"http_server"`
	Clients    ClientsConfig `yaml:"clients"`
	// AppSecret  string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-default:":8000"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout" env-default:"5s"`
	RetriesCount int           `yaml:"retries_count"`
	Insecure     bool          `yaml:"insecure"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}

	return c, nil
}
