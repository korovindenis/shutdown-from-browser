package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	storage "github.com/korovindenis/shutdown-from-browser/v2/internal/adapter/storage/memory"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/config"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/domain/usecase"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/handler"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/middleware"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/middleware/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Exec() error {
	// config
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "config load")
	}

	// logger
	log, err := setupLogger(cfg.Env)
	if err != nil {
		return errors.Wrap(err, "setup logger")
	}

	// bd
	computerStorage, err := storage.New()
	if err != nil {
		return errors.Wrap(err, "bd storage init")
	}

	// domain
	computerUsecase := usecase.NewComputerUsecase(computerStorage, log)
	computerHandler := handler.NewComputerHandler(computerUsecase, cfg, log)

	// run main logic
	go computerUsecase.IsNeedPowerOff(cfg.LogsLevel)

	// http
	gin.SetMode(cfg.HTTPServer.Mode)
	router := gin.Default()
	// middleware
	router.Use(middleware.RequestIdMiddleware())
	router.Use(gin.Recovery())
	router.Use(logger.LoggingMiddleware(log))
	router.LoadHTMLGlob(cfg.HTTPServer.TemplatesPath)
	// Define endpoint
	router.GET("/", computerHandler.MainPageHandler)
	router.Static("/static/", "./web/build/static")
	router.POST("/api/v1/server-power/", computerHandler.SetPowerOffHandler)
	router.GET("/api/v1/get-time-autopoweroff/", computerHandler.GetTimePoHandler)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Failed to listen and serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	return srv.Shutdown(ctx)
}

func setupLogger(env string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if env == "prod" {
		cfg := zap.NewProductionConfig()
		logger, err = cfg.Build()
		if err != nil {
			return nil, errors.Wrap(err, "setupLogger")
		}
	} else {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(encoderConfig)

		consoleOutput := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(encoder, consoleOutput, zapcore.DebugLevel)

		logger = zap.New(core)
	}

	return logger, nil
}
