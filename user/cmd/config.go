package main

import "os"

type Config struct {
	MysqlDsn      string
	Port          string
	MigrationPath string
}

func initConfig() Config {
	return Config{
		MysqlDsn:      os.Getenv("MYSQL_DSN"),
		Port:          os.Getenv("APP_PORT"),
		MigrationPath: os.Getenv("USER_MIGRATION_PATH"),
	}
}
