package main

import "os"

type Config struct {
	MysqlDsn string
	Port     string
}

func initConfig() Config {
	return Config{
		MysqlDsn: os.Getenv("BATTLE_MYSQL_DSN"),
		Port:     os.Getenv("BATTLE_APP_PORT"),
	}
}
