package config

import (
	"fmt"
	"os"
	"strings"
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

func GetEnv() map[string]string {
	var table = make(map[string]string)

	buf, err := os.ReadFile(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	var result = strings.Fields(string(buf))

	for i := range result {
		if strings.Contains(result[i], "=") {
			var res = strings.Split(result[i], "=")
			table[res[0]] = res[1]
		}
	}

	return table
}

func NewDB(table map[string]string) *DB {
	db := &DB{
		Host:     table["POSTGRES_HOST"],
		Port:     table["POSTGRES_PORT"],
		User:     table["POSTGRES_USER"],
		Password: table["POSTGRES_PASSWORD"],
		Name:     table["POSTGRES_DBNAME"],
	}

	return db
}

func NewConfig(table map[string]string) *Config {
	var config = &Config{
		Interval: table["CLI_APP_TIMER_INTERVAL"],
		Workers:  table["CLI_APP_WORKERS_COUNT"],
	}

	return config
}

func NewAppConfig() *AppConfig {
	var table = GetEnv()

	config := &AppConfig{
		DB:     NewDB(table),
		Config: NewConfig(table),
	}

	return config
}
