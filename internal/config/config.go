package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type HTTPConfig struct {
	Addr            string        `yaml:"addr"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type PostgresConfig struct {
	DSN             string        `yaml:"dsn"`
	MaxConns        int32         `yaml:"max_conns"`
	MinConns        int32         `yaml:"min_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load(".env")
	cfg := defaultConfig()

	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			return nil, fmt.Errorf("parse config: %w", err)
		}
	}
	applyEnv(&cfg)
	if cfg.Postgres.DSN == "" {
		return nil, fmt.Errorf("postgres dsn is required")
	}
	return &cfg, nil
}

func defaultConfig() Config {
	return Config{
		HTTP:     HTTPConfig{Addr: ":8080", ShutdownTimeout: 10 * time.Second},
		Postgres: PostgresConfig{MaxConns: 10, MinConns: 1, MaxConnLifetime: time.Hour},
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		cfg.HTTP.Addr = v
	}
	if v := os.Getenv("POSTGRES_DSN"); v != "" {
		cfg.Postgres.DSN = v
	}
	if v := os.Getenv("POSTGRES_MAX_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Postgres.MaxConns = int32(n)
		}
	}
	if v := os.Getenv("POSTGRES_MIN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Postgres.MinConns = int32(n)
		}
	}
	if v := os.Getenv("POSTGRES_MAX_CONN_LIFETIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Postgres.MaxConnLifetime = d
		}
	}
}
