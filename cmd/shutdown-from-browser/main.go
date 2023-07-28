package main

import (
	"log"
	"os"

	"github.com/korovindenis/shutdown-from-browser/v2/internal/app"
)

const (
	EXIT_SUCCESS = iota
	EXIT_CRITICAL
)

func main() {
	if err := app.Exec(); err != nil {
		log.Println(err)
		os.Exit(EXIT_CRITICAL)
	}
	os.Exit(EXIT_SUCCESS)
}
