package config

import (
	"os"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Config struct {
	Interval string
	Workers  string
}

type AppConfig struct {
	DB     *DB
	Config *Config
}

func NewDB() *DB {
	db := &DB{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DBNAME"),
	}

	return db
}

func NewConfig() *Config {
	var config = &Config{
		Interval: os.Getenv("CLI_APP_TIMER_INTERVAL"),
		Workers:  os.Getenv("CLI_APP_WORKERS_COUNT"),
	}

	return config
}

func NewAppConfig() *AppConfig {
	config := &AppConfig{
		DB:     NewDB(),
		Config: NewConfig(),
	}

	return config
}
