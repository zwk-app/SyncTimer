package timer

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type TargetTimer struct {
	current                *Time
	target                 *Time
	remainingTimeInSeconds int
	remainingHours         int
	remainingMinutes       int
	remainingSeconds       int
}

func NewTargetTimer() *TargetTimer {
	t := TargetTimer{
		current:                NewTime(),
		target:                 NewTime(),
		remainingTimeInSeconds: 0,
		remainingHours:         0,
		remainingMinutes:       0,
		remainingSeconds:       0,
	}
	ch := make(chan time.Time)
	t.remainingTimeLoop(ch)
	t.currentTimeLoop(ch)
	return &t
}

func (t *TargetTimer) currentTimeLoop(ch chan time.Time) {
	go func() {
		for {
			current := time.Now()
			t.current.Set(current)
			ch <- current
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func (t *TargetTimer) remainingTimeLoop(ch chan time.Time) {
	go func() {
		for current := range ch {
			d := t.target.time.Sub(current).Round(time.Second)
			h := d / time.Hour
			d -= h * time.Hour
			m := d / time.Minute
			d -= m * time.Minute
			s := d / time.Second
			t.remainingTimeInSeconds = int(s) + (int(m) * 60) + (int(h) * 3600)
			t.remainingHours = int(h)
			t.remainingMinutes = int(m)
			t.remainingSeconds = int(s)
		}
	}()
}

func (t *TargetTimer) GetLocationName() string {
	return t.target.GetLocationName()
}

func (t *TargetTimer) GetCurrentTimeString() string {
	return t.current.GetTimeString()
}
func (t *TargetTimer) GetTargetTimeString() string {
	return t.target.GetTimeString()
}

func (t *TargetTimer) GetRemainingString() string {
	if t.remainingHours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", t.remainingHours, t.remainingMinutes, t.remainingSeconds)
	}
	if t.remainingMinutes > 0 {
		return fmt.Sprintf("%d:%02d", t.remainingMinutes, t.remainingSeconds)
	}
	return fmt.Sprintf("%d", t.remainingSeconds)
}

func (t *TargetTimer) GetRemainingTime() (int, int, int) {
	return t.remainingHours, t.remainingMinutes, t.remainingSeconds
}

func (t *TargetTimer) GetRemainingSeconds() int {
	return t.remainingTimeInSeconds
}

func (t *TargetTimer) SetLocationName(locationName string) bool {
	log.Printf("TargetTimer->SetLocationName '%s'", locationName)
	return t.target.SetLocationName(locationName) && t.current.SetLocationName(locationName)
}

func (t *TargetTimer) SetDelay(h int, m int, s int) bool {
	log.Printf("TargetTimer->SetDelay '%d:%02d:%02d'", h, m, s)
	if CheckTime(h, m, s) {
		d, e := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		if e != nil {
			log.Printf("TargetTimer->SetDelay ERROR: %s", e.Error())
		} else {
			currentTime := time.Now().In(t.target.location)
			t.target.Set(currentTime.Add(d))
			return true
		}
	}
	return false
}

func (t *TargetTimer) SetDelayString(delay string) bool {
	log.Printf("TargetTimer->SetDelayString '%s'", delay)
	delayLen := len(delay)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch delayLen {
	case 1, 2:
		s, se = strconv.Atoi(delay)
		if se != nil {
			hasError = true
		}
	case 3:
		m, me = strconv.Atoi(delay[0:1])
		s, se = strconv.Atoi(delay[1:3])
		if me != nil || se != nil {
			hasError = true
		}
	case 4:
		m, me = strconv.Atoi(delay[0:2])
		s, se = strconv.Atoi(delay[2:4])
		if me != nil || se != nil {
			hasError = true
		}
	case 5:
		h, he = strconv.Atoi(delay[0:1])
		m, me = strconv.Atoi(delay[1:3])
		s, se = strconv.Atoi(delay[3:5])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	case 6:
		h, he = strconv.Atoi(delay[0:2])
		m, me = strconv.Atoi(delay[2:4])
		s, se = strconv.Atoi(delay[4:6])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	}
	if !hasError {
		return t.SetDelay(h, m, s)
	}
	return false
}

func (t *TargetTimer) SetTarget(h int, m int, s int) bool {
	log.Printf("TargetTimer->SetTarget '%d:%02d:%02d'", h, m, s)
	return t.target.SetTime(h, m, s)
}
func (t *TargetTimer) SetTargetString(target string) bool {
	log.Printf("TargetTimer->SetTargetString '%s'", target)
	return t.target.SetTimeString(target)
}

func (t *TargetTimer) Next() {
	locationName := t.target.GetLocationName()
	t.SetLocationName("UTC")
	if t.current.hours < 11 {
		t.SetTarget(11, 15, 0)
	} else if t.current.hours < 14 {
		t.SetTarget(14, 15, 0)
	} else {
		t.SetTarget(11, 15, 0)
	}
	t.SetLocationName(locationName)
}
