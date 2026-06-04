package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
	JWT  JWTConfig
}

type DBConfig struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	DBName                string
	SSLMode               string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

type HTTPConfig struct {
	Port string
}

type JWTConfig struct {
	Secret                    string
	ExpirationDurationMinutes int
}

func Load() Config {
	return Config{
		DB: DBConfig{
			Host:                  getEnv("DB_HOST", "localhost"),
			Port:                  getEnvInt("DB_PORT", 5432),
			Username:              getEnv("DB_USER", "myuser"),
			Password:              getEnv("DB_PASSWORD", "mypassword"),
			DBName:                getEnv("DB_NAME", "mydatabase"),
			SSLMode:               getEnv("DB_SSLMODE", "disable"),
			MaxOpenConnections:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConnections:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnectionMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret:                    getEnv("JWT_SECRET", ""),
			ExpirationDurationMinutes: getEnvInt("JWT_EXPIRATION_DURATION_MINUTES", 24*60), //24 hours
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v, err := time.ParseDuration(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}
