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
	"github.com/korovindenis/shutdown-from-browser/pkg/countdown"
)

/*
    syscall commands
    LINUX_REBOOT_CMD_CAD_OFF         = 0x0
    LINUX_REBOOT_CMD_CAD_ON          = 0x89abcdef
    LINUX_REBOOT_CMD_HALT            = 0xcdef0123
    LINUX_REBOOT_CMD_KEXEC           = 0x45584543
    LINUX_REBOOT_CMD_POWER_OFF       = 0x4321fedc
    LINUX_REBOOT_CMD_RESTART         = 0x1234567
    LINUX_REBOOT_CMD_RESTART2        = 0xa1b2c3d4
    LINUX_REBOOT_CMD_SW_SUSPEND      = 0xd000fce2
    LINUX_REBOOT_MAGIC1              = 0xfee1dead
    LINUX_REBOOT_MAGIC2              = 0x28121969
   Linux Man page info
   LINUX_REBOOT_CMD_CAD_OFF
          (RB_DISABLE_CAD, 0).  CAD is disabled.  This means that the
          CAD keystroke will cause a SIGINT signal to be sent to init
          (process 1), whereupon this process may decide upon a proper
          action (maybe: kill all processes, sync, reboot).
   LINUX_REBOOT_CMD_CAD_ON
          (RB_ENABLE_CAD, 0x89abcdef).  CAD is enabled.  This means that
          the CAD keystroke will immediately cause the action associated
          with LINUX_REBOOT_CMD_RESTART.
   LINUX_REBOOT_CMD_HALT
          (RB_HALT_SYSTEM, 0xcdef0123; since Linux 1.1.76).  The message
          "System halted." is printed, and the system is halted.
          Control is given to the ROM monitor, if there is one.  If not
          preceded by a sync(2), data will be lost.
   LINUX_REBOOT_CMD_KEXEC
          (RB_KEXEC, 0x45584543, since Linux 2.6.13).  Execute a kernel
          that has been loaded earlier with kexec_load(2).  This option
          is available only if the kernel was configured with
          CONFIG_KEXEC.
   LINUX_REBOOT_CMD_POWER_OFF
          (RB_POWER_OFF, 0x4321fedc; since Linux 2.1.30).  The message
          "Power down." is printed, the system is stopped, and all power
          is removed from the system, if possible.  If not preceded by a
          sync(2), data will be lost.
   LINUX_REBOOT_CMD_RESTART
          (RB_AUTOBOOT, 0x1234567).  The message "Restarting system." is
          printed, and a default restart is performed immediately.  If
          not preceded by a sync(2), data will be lost.
   LINUX_REBOOT_CMD_RESTART2
          (0xa1b2c3d4; since Linux 2.1.30).  The message "Restarting
          system with command '%s'" is printed, and a restart (using the
          command string given in arg) is performed immediately.  If not
          preceded by a sync(2), data will be lost.
   LINUX_REBOOT_CMD_SW_SUSPEND
          (RB_SW_SUSPEND, 0xd000fce1; since Linux 2.5.18).  The system
          is suspended (hibernated) to disk.  This option is available
          only if the kernel was configured with CONFIG_HIBERNATION.
*/

type Sfb struct {
	httpServer *http.Server
}

var myServer countdown.Server

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

func powerHandler(w http.ResponseWriter, r *http.Request) {
	var tmpServer countdown.Server
	// validate input
	if err := json.NewDecoder(r.Body).Decode(&tmpServer); (err != nil) || (tmpServer.Mode != "shutdown" && tmpServer.Mode != "reboot") {
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
	myServer = tmpServer
	log.Println("Received:", myServer)

	// send response
	jsonResp, err := json.Marshal(map[string]string{"message": "Server is " + myServer.Mode})
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

func getTimePOHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(map[string]string{"time": "2022-10-26T23:11:45.664Z"}) //myServer.TimeShutDown})
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
