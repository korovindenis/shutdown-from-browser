package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/config"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/middleware"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/middleware/logger"
	"go.uber.org/zap"
)

type usecases interface {
	IsNeedPowerOff(ctx context.Context, logslevel uint8)
}

type handlers interface {
	SetPowerOffHandler(c *gin.Context)
	MainPageHandler(c *gin.Context)
	GetTimePoHandler(c *gin.Context)
}

func Exec(cfg *config.Config, log *zap.Logger, computerUsecase usecases, computerhandler handlers) error {
	// run main logic
	ctxCncl, cancel := context.WithCancel(context.Background())
	defer cancel()
	go computerUsecase.IsNeedPowerOff(ctxCncl, cfg.LogsLevel)

	// prepare http
	gin.SetMode(cfg.HTTPServer.Mode)
	router := gin.Default()

	// middleware
	router.Use(middleware.RequestIdMiddleware())
	router.Use(gin.Recovery())
	router.Use(logger.LoggingMiddleware(log))
	router.LoadHTMLGlob(cfg.HTTPServer.TemplatesPath)

	// Define endpoint
	router.GET("/", computerhandler.MainPageHandler)
	router.Static("/static/", "./web/build/static")
	router.POST("/api/v1/server-power/", computerhandler.SetPowerOffHandler)
	router.GET("/api/v1/get-time-autopoweroff/", computerhandler.GetTimePoHandler)

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
