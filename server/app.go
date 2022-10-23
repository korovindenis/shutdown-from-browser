package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/spf13/viper"
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

// Received from frontend
type PowerMode struct {
	Mode string
	When string
}

func NewSfb() *Sfb {
	// Define the rice box with the frontend static files
	appBox, err := rice.FindBox("../web/build")
	if err != nil {
		log.Fatal(err)
	}

	// Define endpoint
	http.HandleFunc("/", serveAppHandler(appBox))
	http.HandleFunc("/api/v1/server-power/", powerHandler)
	http.HandleFunc("/api/v1/get-time-autopoweroff/", getTimePOHandler)
	// server static files
	http.Handle("/static/", http.FileServer(appBox.HTTPBox()))

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
	// get request
	var pcState PowerMode
	if err := json.NewDecoder(r.Body).Decode(&pcState); (err != nil) || (pcState.Mode != "shutdown" && pcState.Mode != "reboot" || !IsISO8601Date(pcState.When)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Received", pcState.Mode)

	// send response
	response := make(map[string]string)
	response["message"] = "Server is " + pcState.Mode

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error JSON Marshal : %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonResp)

	// bye
	log.Println("Run:", viper.GetString(pcState.Mode))
}

func getTimePOHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "hi"

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Error JSON Marshal : %s", err)
	}
	time.Sleep(8 * time.Second)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(jsonResp)
}

func serveAppHandler(app *rice.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := app.Open("index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, "index.html", time.Time{}, indexFile)
	}
}

func IsISO8601Date(_date string) bool {
	ISO8601DateRegexString := "^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])(?:T|\\s)(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])?(Z)?.[0-9]{3}[A-Z]{1}$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)

	return ISO8601DateRegex.MatchString(_date)
}
