package main

import (
	"os"

	storage "github.com/korovindenis/shutdown-from-browser/v2/internal/adapter/storage/memory"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/app"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/config"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/domain/usecase"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/handler"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/logger"
	"github.com/pkg/errors"
)

const (
	ExitSuccess = iota
	ExitCritical
)

func main() {
	// init config
	cfg, err := config.New()
	if err != nil {
		errors.Wrap(err, "config load")
		os.Exit(ExitCritical)
	}

	// init logger
	log, err := logger.New(cfg.Env)
	if err != nil {
		errors.Wrap(err, "setup logger")
		os.Exit(ExitCritical)
	}

	// init bd
	computerStorage, err := storage.New()
	if err != nil {
		errors.Wrap(err, "bd storage init")
		os.Exit(ExitCritical)
	}

	// init domain
	computerUsecase := usecase.New(computerStorage, log)
	computerHandler := handler.New(computerUsecase, cfg, log)

	if err := app.Exec(cfg, log, computerUsecase, computerHandler); err != nil {
		errors.Wrap(err, "app Exec was failed")
		os.Exit(ExitCritical)
	}
	os.Exit(ExitSuccess)
}
