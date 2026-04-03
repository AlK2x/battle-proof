package main

import "os"

type Config struct {
	MongoDSN string
	Port     string
}

func initConfig() Config {
	return Config{
		MongoDSN: os.Getenv("MONGO_DSN"),
		Port:     os.Getenv("APP_PORT"),
	}
}
