package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// GetEnv возвращает значение переменной окружения или fallback, если она не задана
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// GetEnvInt возвращает значение переменной окружения как int или fallback, если она не задана либо некорректна
func GetEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// GetEnvDuration возвращает значение переменной окружения как time.Duration или fallback, если она не задана либо некорректна
func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

// BuildPostgresURL формирует строку подключения к Postgres из переменных окружения
func BuildPostgresURL() string {
	host := GetEnv("POSTGRES_HOST", "localhost")
	port := GetEnv("POSTGRES_PORT", "5432")
	user := GetEnv("POSTGRES_USER", "postgres")
	password := GetEnv("POSTGRES_PASSWORD", "postgres")
	dbname := GetEnv("POSTGRES_DB", "base_go")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
}
