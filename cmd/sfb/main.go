package main

import (
	"log"
	"os"
	config "github.com/korovindenis/shutdown-from-browser/v1/configs"
	server "github.com/korovindenis/shutdown-from-browser/v1/server"
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
		log.Printf("%s", err.Error())
		os.Exit(1)
	}
	port := viper.GetString("port")
	logslevel := viper.GetUint("logslevel")

	sfb, err := server.NewSfb(logslevel)
	if err != nil {
		log.Printf("%s", err.Error())
		os.Exit(1)
	}
	if err := sfb.Run(port, logslevel); err != nil {
		log.Printf("%s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
