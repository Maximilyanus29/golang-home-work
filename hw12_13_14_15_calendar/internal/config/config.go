package config

import (
	"os"

	"github.com/Maximilyanus29/golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Storage StorageConf
}

type LoggerConf struct {
	Type     string
	Level    string
	FilePath string
}

type ServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type        string
	DsnPostgres string
}

func NewConfig(cfgFile *string) Config {
	if *cfgFile != "" {
		viper.SetConfigFile(*cfgFile)
	} else {
		app.Exit("flag config is required")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		app.Exit(err.Error())
	}

	return Config{
		Logger: LoggerConf{
			Type:     viper.GetString("logger.type"),
			Level:    viper.GetString("logger.log_level"),
			FilePath: viper.GetString("logger.file-path"),
		},
		Server: ServerConf{
			Host: os.Getenv("HOST"),
			Port: os.Getenv("PORT"),
		},
		Storage: StorageConf{
			Type:        viper.GetString("storage.type"),
			DsnPostgres: viper.GetString("storage.dsn-postgres"),
		},
	}
}
