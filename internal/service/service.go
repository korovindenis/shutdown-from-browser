package service

import "time"

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
