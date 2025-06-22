package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type            string
	DsnPostgres     string
	DsnPostgresTest string
}

func NewConfig() Config {
	return Config{
		Logger: LoggerConf{
			Level: viper.GetString("logger.level"),
		},
		Server: ServerConf{
			Host: os.Getenv("HOST"),
			Port: os.Getenv("PORT"),
		},
		Storage: StorageConf{
			Type:            viper.GetString("storage.type"),
			DsnPostgres:     viper.GetString("storage.dsn-postgres"),
			DsnPostgresTest: viper.GetString("storage.dsn-postgres-test"),
		},
	}
}
