package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/korovindenis/shutdown-from-browser/models"
)

var MyServer models.ServerStatus

// PowerHandler
// @Summary      PowerHandler
// @Description  set time for reboot or shutdown
// @Tags         Reboot or shutdown
// @Accept       json
// @Produce      json
// @Param		 input body models.ServerStatus true "format time is RFC3339"
// @Success      200  {object}  models.PoResponse
// @Router       /server-power/ [post]
func PowerHandler(w http.ResponseWriter, r *http.Request) {
	var tmpServer models.ServerStatus
	// validate input
	if err := json.NewDecoder(r.Body).Decode(&tmpServer); (err != nil) || (tmpServer.Mode != "" && tmpServer.Mode != "shutdown" && tmpServer.Mode != "reboot") {
		log.Printf("Error validate server Mode : %s", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// validate time (format and time no more 24h.)
	tomorrow := time.Now().Add(24 * time.Hour).UTC()
	if responseTime, err := time.Parse(time.RFC3339, tmpServer.TimeShutDown); err != nil || responseTime.After(tomorrow) {
		log.Printf("Error validate timestamp : %s", err)

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if resived too many requests
	if tmpServer.TimeShutDown == MyServer.TimeShutDown {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	MyServer = tmpServer
	log.Printf("Received: %+v", MyServer)

	// send response
	jsonResp, err := json.Marshal(models.PoResponse{Message: "Server is " + MyServer.Mode + " on the " + MyServer.TimeShutDown})
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

// GetTimePOHandler
// @Summary      GetTimePOHandler
// @Description  get the auto power off time
// @Tags         Get time
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ServerStatus
// @Router       /get-time-autopoweroff/ [get]
func GetTimePOHandler(w http.ResponseWriter, r *http.Request) {
	var res models.ServerStatus
	if MyServer.Mode != "" {
		res = MyServer
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
