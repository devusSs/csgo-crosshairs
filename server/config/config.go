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

	RedisHost     string `json:"redis_host"`
	RedisPort     int    `json:"redis_port"`
	RedisPassword string `json:"redis_password"`

	MinioHost     string `json:"minio_host"`
	MinioPort     int    `json:"minio_port"`
	MinioDomain   string `json:"minio_domain"`
	MinioUser     string `json:"minio_user"`
	MinioPassword string `json:"minio_password"`

	APIHost string `json:"api_host"`
	APIPort int    `json:"api_port"`
	Domain  string `json:"domain"`

	SecretSessionsKey string `json:"secret_sessions_key"`

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

	if c.RedisHost == "" {
		return errors.New("missing key: redis_host")
	}

	if c.RedisPort == 0 {
		return errors.New("missing key: redis_port")
	}

	if c.RedisPassword == "" {
		return errors.New("missing key: redis_password")
	}

	if c.MinioHost == "" {
		if c.MinioDomain == "" {
			return errors.New("missing key: minio_host")
		}
	}

	if c.MinioPort == 0 {
		if c.MinioDomain == "" {
			return errors.New("missing key: minio_port")
		}
	}

	if c.MinioUser == "" {
		return errors.New("missing key: minio_user")
	}

	if c.MinioPassword == "" {
		return errors.New("missing key: minio_password")
	}

	if c.APIPort == 0 {
		return errors.New("missing key: api_port")
	}

	if c.Domain == "" {
		return errors.New("missing key: domain")
	}

	if c.SecretSessionsKey == "" {
		return errors.New("missing key: secret_sessions_key")
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
