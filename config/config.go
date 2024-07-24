package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Environment string

// each environment has a dedicated config
// the base config is loaded as a default with the environment specific config loaded on top
const (
	Local Environment = "local" // contains overrides for local development
	Test  Environment = "test"  // contains overrides for the test environment
	Acc   Environment = "acc"   // contains overrides for the test environment
	Prod  Environment = "prod"  // contains overrides for the production environment
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl-mode"`
		Name     string `yaml:"name"`
		Port     int    `yaml:"port"`
	} `yaml:"database"`
}

func Load(env Environment) (*Config, error) {
	f, err := os.Open(fmt.Sprintf("%s.yml", env))
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode yaml config file: %w", err)
	}

	return &cfg, nil
}

// DBConnString returns the connection string in the format: "host=host port=port sslmode=mode dbname=dbname user=user password=pass"
func (c *Config) DBConnString() string {
	return c.DBConnStringWithDBName(c.Database.Name)
}

// DBConnStringWithDBName returns the connection string in the format:
// "host=host port=port sslmode=mode dbname=name user=user password=pass"
func (c *Config) DBConnStringWithDBName(name string) string {
	connection := fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.SSLMode,
		name,
	)

	if c.Database.User != "" {
		connection += fmt.Sprintf(" user=%s", c.Database.User)
	}

	if c.Database.Password != "" {
		connection += fmt.Sprintf(" password=%s", c.Database.Password)
	}

	return connection
}
