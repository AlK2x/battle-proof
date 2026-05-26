package config

import (
	"os"
	"strconv"
)

type Config struct {
	MysqlDsn  string
	Port      string
	KafkaAddr string
	Debug     bool
}

func Init() Config {
	debug := os.Getenv("APP_DEBUG")
	isDebug, _ := strconv.ParseBool(debug)
	return Config{
		MysqlDsn: os.Getenv("MYSQL_DSN"),
		Port:     os.Getenv("APP_PORT"),
		Debug:    isDebug,
	}
}
