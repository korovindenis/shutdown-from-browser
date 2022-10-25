package countdown

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	Mode         string
	TimeShutDown string
}

type countdown struct {
	t int
	d int
	h int
	m int
	s int
}

func New(s *Server) {
	for s.Mode == "" {
		time.Sleep(time.Second * 5)
	}
	for range time.Tick(1 * time.Second) {
		v, _ := time.Parse(time.RFC3339, s.TimeShutDown)
		timeRemaining := getTimeRemaining(v)

		if timeRemaining.t <= 0 {
			// bye
			log.Println("Run:", viper.GetString(s.Mode))
		}

		log.Printf("Time for %s - Days: %d Hours: %d Minutes: %d Seconds: %d\n", s.Mode, timeRemaining.d, timeRemaining.h, timeRemaining.m, timeRemaining.s)
	}

}

func getTimeRemaining(t time.Time) countdown {
	currentTime := time.Now()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	days := int(total / (60 * 60 * 24))
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return countdown{
		t: total,
		d: days,
		h: hours,
		m: minutes,
		s: seconds,
	}
}
