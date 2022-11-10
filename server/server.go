package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/korovindenis/shutdown-from-browser/api"
	"github.com/korovindenis/shutdown-from-browser/pkg/countdown"
	"github.com/korovindenis/shutdown-from-browser/pkg/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Sfb struct {
	HttpServer *http.Server
	WebFolder  *rice.Box
}

// Define the rice box with the frontend static files
func (s *Sfb) FindBox() (res bool, err error) {
	if s.WebFolder, err = rice.FindBox("../web/build"); err != nil {
		return res, err
	}
	return true, nil
}

func (s *Sfb) Run(port string) error {
	s.HttpServer = &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("Server starting at port", port)

		if err := s.HttpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.HttpServer.Shutdown(ctx)
}

func NewSfb() *Sfb {
	newApp := Sfb{}
	newApp.FindBox()

	// Define endpoint
	http.HandleFunc("/", serveAppHandler(newApp.WebFolder))
	http.HandleFunc("/api/v1/server-power/", handler.PowerHandler)
	http.HandleFunc("/api/v1/get-time-autopoweroff/", handler.GetTimePOHandler)
	// Server static files
	http.Handle("/static/", http.FileServer(newApp.WebFolder.HTTPBox()))
	// Swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	go countdown.New(&handler.MyServer)

	return &newApp
}

func serveAppHandler(app *rice.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := app.Open("index.html")
		if err != nil {
			log.Printf("Error open index.html : %s", err)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
	}
}
