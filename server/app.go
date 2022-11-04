package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	rice "github.com/GeertJohan/go.rice"
	_ "github.com/korovindenis/shutdown-from-browser/api"
	"github.com/korovindenis/shutdown-from-browser/models"
	"github.com/korovindenis/shutdown-from-browser/pkg/countdown"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Sfb struct {
	httpServer *http.Server
}

var myServer models.Server

func NewSfb() *Sfb {
	// Define the rice box with the frontend static files
	appBox, err := rice.FindBox("../web/build")
	if err != nil {
		log.Fatal(err)
	}

	// // Define endpoint
	http.HandleFunc("/", serveAppHandler(appBox))
	http.HandleFunc("/api/v1/server-power/", powerHandler)
	http.HandleFunc("/api/v1/get-time-autopoweroff/", getTimePOHandler)
	// server static files
	http.Handle("/static/", http.FileServer(appBox.HTTPBox()))
	// swagger
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	go countdown.New(&myServer)

	return &Sfb{}
}

func (s *Sfb) Run(port string) error {

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("Server starting at port", port)

		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}

// PowerHandler
// @Summary      PowerHandler
// @Description  set time for reboot or shutdown
// @Tags         Reboot or shutdown
// @Accept       json
// @Produce      json
// @Param		 input body models.Server true "format time is RFC3339"
// @Success      200  {object}  models.PoResponse
// @Router       /server-power/ [post]
func powerHandler(w http.ResponseWriter, r *http.Request) {
	var tmpServer models.Server
	// validate input
	if err := json.NewDecoder(r.Body).Decode(&tmpServer); (err != nil) || (tmpServer.Mode != "" && tmpServer.Mode != "shutdown" && tmpServer.Mode != "reboot") {
		log.Printf("Error validate server Mode : %s", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// validate timestamp
	if _, err := time.Parse(time.RFC3339, tmpServer.TimeShutDown); err != nil {
		log.Printf("Error validate timestamp : %s", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if resived too many requests
	if tmpServer.TimeShutDown == myServer.TimeShutDown {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// send response
	jsonResp, err := json.Marshal(models.PoResponse{Message: "Server is " + myServer.Mode + " on the " + myServer.TimeShutDown})
	if err != nil {
		log.Printf("Error JSON Marshal : %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	myServer = tmpServer
	log.Printf("Received: %+v", myServer)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Error return json : %s", err)
	}
}

// GetTimePOHandler
// @Summary      GetTimePOHandler
// @Description  get the auto power off time
// @Tags         Get time
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Server
// @Router       /get-time-autopoweroff/ [get]
func getTimePOHandler(w http.ResponseWriter, r *http.Request) {
	var res models.Server
	if myServer.Mode != "" {
		res = myServer
	}
	jsonResp, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error JSON Marshal : %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Error return json : %s", err)
	}
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
