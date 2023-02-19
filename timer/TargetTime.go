package timer

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const TargetTimeDefaultTextLabel = "Target"

type TargetTime struct {
	current    *Time
	target     *Time
	textLabel  string
	alertSound string
	remaining  struct {
		timeInSeconds int
		hours         int
		minutes       int
		seconds       int
	}
	last struct {
		error error
	}
}

func NewTargetTimer() *TargetTime {
	t := TargetTime{}
	t.current = NewTime()
	t.target = NewTime()
	t.textLabel = ""
	t.alertSound = ""
	t.remaining.timeInSeconds = 0
	t.remaining.hours = 0
	t.remaining.minutes = 0
	t.remaining.seconds = 0
	ch := make(chan time.Time)
	t.remainingTimeLoop(ch)
	t.currentTimeLoop(ch)
	return &t
}

func (t *TargetTime) Error() error {
	lastError := t.last.error
	t.last.error = nil
	return lastError
}

func (t *TargetTime) currentTimeLoop(ch chan time.Time) {
	go func() {
		for {
			current := time.Now()
			t.current.Set(current)
			ch <- current
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func (t *TargetTime) remainingTimeLoop(ch chan time.Time) {
	go func() {
		for current := range ch {
			d := t.target.time.Sub(current).Round(time.Second)
			h := d / time.Hour
			d -= h * time.Hour
			m := d / time.Minute
			d -= m * time.Minute
			s := d / time.Second
			t.remaining.timeInSeconds = int(s) + (int(m) * 60) + (int(h) * 3600)
			t.remaining.hours = int(h)
			t.remaining.minutes = int(m)
			t.remaining.seconds = int(s)
		}
	}()
}

func (t *TargetTime) TextLabel() string {
	if t.textLabel == "" {
		return TargetTimeDefaultTextLabel
	}
	return t.textLabel
}

func (t *TargetTime) AlertSound() string {
	return t.alertSound
}

func (t *TargetTime) LocationName() string {
	return t.target.LocationName()
}

func (t *TargetTime) CurrentTimeString() string {
	return t.current.TimeString()
}

func (t *TargetTime) CurrentTextString() string {
	return t.current.TextString()
}

func (t *TargetTime) TargetTimeString() string {
	return t.target.TimeString()
}

func (t *TargetTime) TargetTextString() string {
	return t.target.TextString()
}

func (t *TargetTime) RemainingString() string {
	if t.remaining.hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", t.remaining.hours, t.remaining.minutes, t.remaining.seconds)
	}
	if t.remaining.minutes > 0 {
		return fmt.Sprintf("%d:%02d", t.remaining.minutes, t.remaining.seconds)
	}
	return fmt.Sprintf("%d", t.remaining.seconds)
}

func (t *TargetTime) RemainingTime() (int, int, int) {
	return t.remaining.hours, t.remaining.minutes, t.remaining.seconds
}

func (t *TargetTime) RemainingSeconds() int {
	return t.remaining.timeInSeconds
}

func (t *TargetTime) SetTextLabel(textLabel string) *TargetTime {
	t.textLabel = textLabel
	return t
}

func (t *TargetTime) SetAlertSound(alertSound string) *TargetTime {
	t.alertSound = alertSound
	return t
}
func (t *TargetTime) SetLocationName(locationName string) *TargetTime {
	log.Printf("TargetTime->SetLocationName '%s'", locationName)
	targetLocation := t.target.LocationName()
	currentLocation := t.current.LocationName()
	if targetLocation != locationName || currentLocation != locationName {
		t.target.SetLocationName(locationName)
		t.current.SetLocationName(locationName)
	}
	return t
}

func (t *TargetTime) SetDelay(h int, m int, s int) *TargetTime {
	log.Printf("TargetTime->SetDelay '%d:%02d:%02d'", h, m, s)
	if !CheckTime(h, m, s) {
		t.last.error = fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
		log.Printf("TargetTime->SetDelay error: %s", t.last.error.Error())
	} else {
		d, e := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		if e != nil {
			t.last.error = e
			log.Printf("TargetTime->SetDelay error: %s", t.last.error.Error())
		} else {
			currentTime := time.Now().In(t.target.location)
			t.target.Set(currentTime.Add(d))
			if t.last.error != nil {
				log.Printf("TargetTime->SetDelay error: %s", t.last.error.Error())
			}
		}
	}
	return t
}

func (t *TargetTime) SetDelayString(delayString string) *TargetTime {
	if delayString == "" {
		return t
	}
	log.Printf("TargetTime->SetDelayString '%s'", delayString)
	delayLen := len(delayString)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch delayLen {
	case 1, 2:
		s, se = strconv.Atoi(delayString)
		if se != nil {
			hasError = true
		}
	case 3:
		m, me = strconv.Atoi(delayString[0:1])
		s, se = strconv.Atoi(delayString[1:3])
		if me != nil || se != nil {
			hasError = true
		}
	case 4:
		m, me = strconv.Atoi(delayString[0:2])
		s, se = strconv.Atoi(delayString[2:4])
		if me != nil || se != nil {
			hasError = true
		}
	case 5:
		h, he = strconv.Atoi(delayString[0:1])
		m, me = strconv.Atoi(delayString[1:3])
		s, se = strconv.Atoi(delayString[3:5])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	case 6:
		h, he = strconv.Atoi(delayString[0:2])
		m, me = strconv.Atoi(delayString[2:4])
		s, se = strconv.Atoi(delayString[4:6])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	}
	if hasError {
		t.last.error = fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
		log.Printf("TargetTime->SetDelayString error: %s", t.last.error.Error())
	} else {
		t.SetDelay(h, m, s)
		if t.last.error != nil {
			log.Printf("TargetTime->SetDelayString error: %s", t.last.error.Error())
		}
	}
	return t
}

func (t *TargetTime) SetTarget(v *Time) *TargetTime {
	log.Printf("TargetTime->SetTarget '%s'", v.TimeString())
	t.target = v
	return t
}

func (t *TargetTime) SetTargetTime(h int, m int, s int) *TargetTime {
	log.Printf("TargetTime->SetTargetTime '%d:%02d:%02d'", h, m, s)
	t.target.SetTime(h, m, s)
	if t.last.error != nil {
		log.Printf("TargetTime->SetTargetTime error: %s", t.last.error.Error())
	}
	return t
}
func (t *TargetTime) SetTargetString(targetString string) *TargetTime {
	if targetString == "" {
		return t
	}
	log.Printf("TargetTime->SetTargetString '%s'", targetString)
	t.target.SetTimeString(targetString)
	if t.last.error != nil {
		log.Printf("TargetTime->SetTargetString error: %s", t.last.error.Error())
	}
	return t
}
