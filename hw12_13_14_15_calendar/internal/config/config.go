package config

import (
	"log"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
}

type LoggerConf struct {
	Level string
}

func NewConfig(cfgFile string) Config {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.SetConfigType("toml")            // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path.Dir(cfgFile)) // path to look for the config file in
	err := viper.ReadInConfig()            // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		log.Printf("fatal error config file: %s", err)
	}

	return Config{
		Logger: LoggerConf{
			Level: viper.GetString("logger.level"),
		},
	}
}
