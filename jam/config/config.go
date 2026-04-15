package config

import "os"

type Config struct {
	MysqlDsn  string
	Port      string
	KafkaAddr string
}

func Init() Config {
	return Config{
		MysqlDsn: os.Getenv("MYSQL_DSN"),
		Port:     os.Getenv("APP_PORT"),
	}
}
