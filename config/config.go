package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
}

func LoadConfig() (config Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Cannot read config file", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Unable to decode into struct", err)
	}
	return
}
