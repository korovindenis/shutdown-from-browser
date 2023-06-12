package transport

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/korovindenis/shutdown-from-browser/v1/api"
	"github.com/korovindenis/shutdown-from-browser/v1/internal/service"
	"github.com/korovindenis/shutdown-from-browser/v1/internal/transport/rest/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Sfb struct {
	HttpServer *http.Server `json:"httpserver"`
	WebFolder  *rice.Box    `json:"webfolder"`
}

// Define the rice box with the frontend static files
func (s *Sfb) FindBox() (res bool, err error) {
	if s.WebFolder, err = rice.FindBox("../../web/build"); err != nil {
		return res, err
	}
	return true, nil
}

func (s *Sfb) Run(port uint32, logslevel uint) error {
	s.HttpServer = &http.Server{
		Addr:           ":" + string(port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if logslevel > 0 {
			log.Println("Server starting at port", port)
		}

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

func NewSfb(logslevel uint) (*Sfb, error) {
	newApp := Sfb{}
	if _, err := newApp.FindBox(); err != nil {
		return nil, err
	}

	// Define endpoint
	http.HandleFunc("/", serveAppHandler(newApp.WebFolder))
	http.HandleFunc("/api/v1/server-power/", handler.PowerHandler)
	http.HandleFunc("/api/v1/get-time-autopoweroff/", handler.GetTimePOHandler)

	// Server static files
	http.Handle("/static/", http.FileServer(newApp.WebFolder.HTTPBox()))

	// Swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	go service.New(&handler.MyServer, logslevel)

	return &newApp, nil
}

func Exec(port uint32, logslevel uint) error {
	sfb, err := NewSfb(logslevel)
	if err != nil {
		return err
	}

	err = sfb.Run(port, logslevel)
	if err != nil {
		return err
	}

	return nil
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
