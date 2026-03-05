package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env       string          `yaml:"env"`
	Server    ServerConfig    `yaml:"server"`
	HTTP      HTTPConfig      `yaml:"http"`
	Storage   StorageConfig   `yaml:"storage"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Shortener ShortenerConfig `yaml:"shortener"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type HTTPConfig struct {
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

type StorageConfig struct {
	Memory   MemoryConfig   `yaml:"memory"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type MemoryConfig struct {
	MaxSize int `yaml:"max_size"`
}

type PostgresConfig struct {
	MaxConnections int32 `yaml:"max_connections"`
	MinConnections int32 `yaml:"min_connections"`
	DSN            string
}

type RateLimitConfig struct {
	RequestsPerSecond int `yaml:"requests_per_second"`
	Burst             int `yaml:"burst"`
}

type ShortenerConfig struct {
	ShortLength int    `yaml:"short_length"`
	Alphabet    string `yaml:"alphabet"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	cfg.Storage.Postgres.DSN = configPostgresDSN()

	return cfg, nil
}

func configPostgresDSN() string {
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	db := getEnv("POSTGRES_DB", "link_shortener")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, db,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
