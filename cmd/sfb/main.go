package main

import (
	"log"

	config "github.com/korovindenis/shutdown-from-browser/configs"
	"github.com/korovindenis/shutdown-from-browser/server"
	"github.com/spf13/viper"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	sfb := server.NewSfb()

	if err := sfb.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
