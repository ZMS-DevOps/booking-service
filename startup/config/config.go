package config

import "os"

type Config struct {
	Port              string
	GrpcPort          string
	BookingDBHost     string
	BookingDBPort     string
	BookingDBUsername string
	BookingDBPassword string
}

func NewConfig() *Config {
	return &Config{
		Port:              os.Getenv("SERVICE_PORT"),
		BookingDBHost:     os.Getenv("DB_HOST"),
		BookingDBPort:     os.Getenv("DB_PORT"),
		BookingDBUsername: os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		BookingDBPassword: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
	}
}
