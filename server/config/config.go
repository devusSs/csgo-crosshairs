package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Config struct {
	PostgresHost     string `json:"postgres_host"`
	PostgresPort     int    `json:"postgres_port"`
	PostgresUser     string `json:"postgres_user"`
	PostgresPassword string `json:"postgres_password"`
	PostgresDB       string `json:"postgres_database"`

	APIHost       string `json:"api_host"`
	APIPort       int    `json:"api_port"`
	BackendDomain string `json:"backend_domain"`

	SecretSessionsKey string `json:"secret_sessions_key"`
	AdminKey          string `json:"admin_key"`

	EmailFrom string `json:"email_from"`
	SMTPHost  string `json:"smtp_host"`
	SMTPPass  string `json:"smtp_pass"`
	SMTPPort  int    `json:"smtp_port"`
	SMTPUser  string `json:"smtp_user"`
}

func LoadConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	jsonConfig, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := json.Unmarshal(jsonConfig, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) CheckConfig() error {
	if c.PostgresHost == "" {
		return errors.New("missing key: postgres_host")
	}

	if c.PostgresPort == 0 {
		return errors.New("missing key: postgres_port")
	}

	if c.PostgresUser == "" {
		return errors.New("missing key: postgres_user")
	}

	if c.PostgresPassword == "" {
		return errors.New("missing key: postgres_password")
	}

	if c.PostgresDB == "" {
		return errors.New("missing key: postgres_database")
	}

	if c.APIHost == "" {
		return errors.New("missing key: api_host")
	}

	if c.APIPort == 0 {
		return errors.New("missing key: api_port")
	}

	if c.BackendDomain == "" {
		return errors.New("missing key: backend_domain")
	}

	if c.SecretSessionsKey == "" {
		return errors.New("missing key: secret_sessions_key")
	}

	if c.AdminKey == "" {
		return errors.New("missing key: admin_key")
	}

	if c.EmailFrom == "" {
		return errors.New("missing key: email_from")
	}

	if c.SMTPHost == "" {
		return errors.New("missing key: smtp_host")
	}

	if c.SMTPPass == "" {
		return errors.New("missing key: smtp_pass")
	}

	if c.SMTPPort == 0 {
		return errors.New("missing key: smtp_port")
	}

	if c.SMTPUser == "" {
		return errors.New("missing key: smtp_user")
	}

	return nil
}
