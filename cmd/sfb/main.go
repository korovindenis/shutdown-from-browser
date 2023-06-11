package main

import (
	"log"
	"os"

	config "github.com/korovindenis/shutdown-from-browser/v1/configs"
	transport "github.com/korovindenis/shutdown-from-browser/v1/internal/transport"
	"github.com/spf13/viper"
)

const (
	EXIT_SUCCESS = 0
	EXIT_ERROR   = 1
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
		os.Exit(EXIT_ERROR)
	}
	port := viper.GetString("port")
	logslevel := viper.GetUint("logslevel")

	sfb, err := transport.NewSfb(logslevel)
	if err != nil {
		log.Printf("%s", err.Error())
		os.Exit(EXIT_ERROR)
	}
	if err := sfb.Run(port, logslevel); err != nil {
		log.Printf("%s", err.Error())
		os.Exit(EXIT_ERROR)
	}
	os.Exit(EXIT_SUCCESS)
}
