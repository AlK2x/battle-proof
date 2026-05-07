package main

import (
	"os"
	"time"
)

type Config struct {
	MysqlDsn      string
	Port          string
	MigrationPath string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

func initConfig() Config {
	return Config{
		MysqlDsn:      os.Getenv("BATTLE_MYSQL_DSN"),
		Port:          os.Getenv("BATTLE_APP_PORT"),
		MigrationPath: os.Getenv("BATTLE_MIGRATION_PATH"),
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  5 * time.Second,
	}
}
