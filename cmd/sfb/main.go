package main

import (
	"log"
	"os"

	config "github.com/korovindenis/shutdown-from-browser/v1/configs"
	transport "github.com/korovindenis/shutdown-from-browser/v1/internal/transport"
	"github.com/spf13/viper"
)

const (
	EXIT_SUCCESS = iota
	EXIT_ERROR
)

// default values
var port uint32 = 8000
var logslevel uint = 1

// read config/config.yml
func init() {
	if err := config.Init(); err != nil {
		log.Printf("%s", err.Error())
	} else {
		port = viper.GetUint32("port")
		logslevel = viper.GetUint("logslevel")
	}
}

// @title           Shutdown from browser
// @version         0.1
// @description     Linux service for shutdown PC from the browser (Go, React)

// @host      localhost:8000
// @BasePath  /api/v1

// @contact.name   korovindenis
// @contact.url    https://github.com/korovindenis
func main() {
	if err := transport.Exec(port, logslevel); err != nil {
		log.Printf("%s", err.Error())
		os.Exit(EXIT_ERROR)
	}
	os.Exit(EXIT_SUCCESS)
}
