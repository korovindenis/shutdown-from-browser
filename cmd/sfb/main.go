package main

import (
	"log"

	config "github.com/korovindenis/shutdown-from-browser/configs"
	"github.com/korovindenis/shutdown-from-browser/server"
	"github.com/spf13/viper"
)

// @title           Shutdown from browser
// @version         0.1
// @description     Linux service for shutdown PC from the browser (Go, React)

// @host      localhost:8000
// @BasePath  /api/v1

// @contact.name   korovindenis
// @contact.url    https://github.com/korovindenis
func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	sfb := server.NewSfb()

	if err := sfb.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
