package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

const homePath = "/usr/bin"

func Init() error {
	if _, err := os.Stat("./configs"); !os.IsNotExist(err) {
		viper.AddConfigPath("./configs")
	} else if _, err := os.Stat(homePath + "/sfb_configs"); !os.IsNotExist(err) {
		viper.AddConfigPath(homePath + "/sfb_configs")
	} else {
		log.Fatalf("Configs folder not found")
	}
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
