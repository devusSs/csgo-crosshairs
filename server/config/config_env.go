package config

import (
	"os"
	"strconv"
	"strings"
)

func LoadEnvConfig() (*Config, error) {
	postgresPort, err := getEnvInt("postgres_port")
	if err != nil {
		return nil, err
	}

	redisPort, err := getEnvInt("redis_port")
	if err != nil {
		return nil, err
	}

	minioPort, err := getEnvInt("minio_port")
	if err != nil {
		return nil, err
	}

	apiPort, err := getEnvInt("api_port")
	if err != nil {
		return nil, err
	}

	smtpPort, err := getEnvInt("smtp_port")
	if err != nil {
		return nil, err
	}

	usingReverseProxy, err := getEnvBool("using_reverse_proxy")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		PostgresHost:     getEnvString("postgres_host"),
		PostgresPort:     postgresPort,
		PostgresUser:     getEnvString("postgres_user"),
		PostgresPassword: getEnvString("postgres_password"),
		PostgresDB:       getEnvString("postgres_db"),

		RedisHost:     getEnvString("redis_host"),
		RedisPort:     redisPort,
		RedisPassword: getEnvString("redis_password"),

		MinioHost:     getEnvString("minio_host"),
		MinioPort:     minioPort,
		MinioDomain:   getEnvString("minio_domain"),
		MinioUser:     getEnvString("minio_root_user"),
		MinioPassword: getEnvString("minio_root_password"),

		APIHost:           getEnvString("api_host"),
		APIPort:           apiPort,
		Domain:            getEnvString("domain"),
		SecretSessionsKey: getEnvString("secret_sessions_key"),

		EmailFrom: getEnvString("email_from"),
		SMTPHost:  getEnvString("smtp_host"),
		SMTPPass:  getEnvString("smtp_pass"),
		SMTPPort:  smtpPort,
		SMTPUser:  getEnvString("smtp_user"),

		UsingReverseProxy: usingReverseProxy,
	}

	return cfg, nil
}

func getEnvString(name string) string {
	return os.Getenv(strings.ToUpper(name))
}

func getEnvInt(name string) (int, error) {
	valueStr := os.Getenv(strings.ToUpper(name))
	return strconv.Atoi(valueStr)
}

func getEnvBool(name string) (bool, error) {
	valueStr := os.Getenv(strings.ToUpper(name))
	return strconv.ParseBool(valueStr)
}
