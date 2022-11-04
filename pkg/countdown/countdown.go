package countdown

import (
	"log"
	"time"

	"github.com/korovindenis/shutdown-from-browser/models"
	"github.com/spf13/viper"
)

type countdown struct {
	t int
	h int
	m int
	s int
}

func New(s *models.Server) {
	for {
		if s.Mode == "" {
			time.Sleep(time.Second * 5)
		} else {
			time.Sleep(1 * time.Second)
			v, _ := time.Parse(time.RFC3339, s.TimeShutDown)
			timeRemaining := getTimeRemaining(v)

			if s.Mode != "" {
				if timeRemaining.t <= 0 {
					// bye
					log.Println("Run:", viper.GetString(s.Mode))
				}
				log.Printf("Time for %s - %d:%d:%d\n", s.Mode, timeRemaining.h, timeRemaining.m, timeRemaining.s)
			}
		}
	}
}

func getTimeRemaining(t time.Time) countdown {
	currentTime := time.Now().UTC()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return countdown{
		t: total,
		h: hours,
		m: minutes,
		s: seconds,
	}
}
