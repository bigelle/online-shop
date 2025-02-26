package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseDSN string
}

func New() *Config {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_USER"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_DBNAME"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_SSLMODE"))
	return &Config{
		DatabaseDSN: dsn,
	}
}
