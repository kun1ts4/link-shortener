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
	Type     string
	Memory   MemoryConfig   `yaml:"memory"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type MemoryConfig struct {
	MaxSize int `yaml:"max_size"`
}

type PostgresConfig struct {
	MaxConnections int32  `yaml:"max_connections"`
	MinConnections int32  `yaml:"min_connections"`
	RetryCount     int    `yaml:"retry_count"`
	RetryDelay     string `yaml:"retry_delay"`
	DSN            string
}

type RateLimitConfig struct {
	RequestsPerSecond int `yaml:"requests_per_second"`
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

	cfg.Storage.Type = getEnv("STORAGE_TYPE", "memory")
	cfg.Storage.Postgres.DSN = configPostgresDSN()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Env == "" {
		return fmt.Errorf("env is required")
	}
	if c.Server.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}
	if c.HTTP.ReadTimeout == "" {
		return fmt.Errorf("http.read_timeout is required")
	}
	if c.HTTP.WriteTimeout == "" {
		return fmt.Errorf("http.write_timeout is required")
	}
	if c.HTTP.IdleTimeout == "" {
		return fmt.Errorf("http.idle_timeout is required")
	}
	if c.Storage.Type == "" {
		return fmt.Errorf("storage.type is required")
	}
	if c.Storage.Type == "memory" && c.Storage.Memory.MaxSize <= 0 {
		return fmt.Errorf("storage.memory.max_size must be greater than 0")
	}
	if c.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("rate_limit.requests_per_second must be greater than 0")
	}
	if c.Shortener.ShortLength <= 0 {
		return fmt.Errorf("shortener.short_length must be greater than 0")
	}
	if c.Shortener.Alphabet == "" {
		return fmt.Errorf("shortener.alphabet is required")
	}
	return nil
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
