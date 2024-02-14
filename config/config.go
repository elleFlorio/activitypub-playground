package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var Domain string
var Port string

func ReadConfig() {
	configFile := os.Getenv("CONFIGFILE")
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	Domain = viper.GetString("domain")
	Port = viper.GetString("port")
}
