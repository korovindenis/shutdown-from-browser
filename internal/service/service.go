package service

import (
	"log"
	"syscall"
	"time"
)

type Status struct {
	Mode         string
	TimeShutDown string
}

type countdown struct {
	t int
	h int
	m int
	s int
}

func New(s *Status, logslevel uint) {
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
					if logslevel > 0 {
						log.Println("Run:", s.Mode)
					}
					callMode := syscall.LINUX_REBOOT_CMD_POWER_OFF
					if s.Mode == "reboot" {
						callMode = syscall.LINUX_REBOOT_CMD_RESTART
					}
					err := syscall.Reboot(callMode)
					if err != nil {
						log.Fatalf("%s", err)
					}
				} else if logslevel > 1 {
					log.Printf("Time for %s - %d:%d:%d\n", s.Mode, timeRemaining.h, timeRemaining.m, timeRemaining.s)
				}
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
