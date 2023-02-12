package ttm

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type Time struct {
	location     *time.Location
	locationName string
	time         time.Time
	hours        int
	minutes      int
	seconds      int
}

const ttmLocalLocationString = "Local Time"

var CurrentTargetTime *Time

func NewTime() *Time {
	return &Time{
		location:     time.Local,
		locationName: ttmLocalLocationString,
		time:         time.Now().In(time.Local),
		hours:        time.Now().In(time.Local).Hour(),
		minutes:      time.Now().In(time.Local).Minute(),
		seconds:      time.Now().In(time.Local).Second(),
	}
}

func (t *Time) GetLocation() *time.Location {
	return t.location
}

func (t *Time) GetLocationName() string {
	return t.locationName
}

func (t *Time) GetTime() time.Time {
	return t.time
}

func GetTimeString(h int, m int, s int) string {
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

func (t *Time) GetTimeString() string {
	return GetTimeString(t.hours, t.minutes, t.seconds)
}

func (t *Time) SetLocation(location string) bool {
	log.Printf("Time->SetLocation '%s'", location)
	if t.locationName != location {
		switch location {
		case "UTC":
			t.location = time.UTC
			t.locationName = "UTC"
		default:
			t.location = time.Local
			t.locationName = ttmLocalLocationString
		}
		t.time = t.time.In(t.location)
		t.hours = t.time.Hour()
		t.minutes = t.time.Minute()
		t.seconds = t.time.Second()
		return true
	}
	return false
}

func CheckTime(h int, m int, s int) bool {
	if h >= 0 && h <= 23 && m >= 0 && m <= 59 && s >= 0 && s <= 59 {
		return true
	}
	return false
}

func (t *Time) SetTime(h int, m int, s int) bool {
	log.Printf("Time->SetTime '%d:%02d:%02d'", h, m, s)
	if CheckTime(h, m, s) {
		currentTime := time.Now().In(t.location)
		tmpTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), h, m, s, 0, t.location)
		if tmpTime.Before(currentTime) {
			tmpTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, h, m, s, 0, t.location)
		}
		t.time = tmpTime.In(t.location)
		t.hours = t.time.Hour()
		t.minutes = t.time.Minute()
		t.seconds = t.time.Second()
		return true
	}
	return false
}

func (t *Time) SetTimeString(target string) bool {
	log.Printf("Time->SetTimeString '%s'", target)
	targetLen := len(target)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch targetLen {
	case 1, 2:
		h, he = strconv.Atoi(target)
		if he != nil {
			hasError = true
		}
	case 4:
		h, he = strconv.Atoi(target[0:2])
		m, me = strconv.Atoi(target[2:4])
		if he != nil || me != nil {
			hasError = true
		}
	case 6:
		h, he = strconv.Atoi(target[0:2])
		m, me = strconv.Atoi(target[2:4])
		s, se = strconv.Atoi(target[4:6])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	}
	if !hasError {
		return t.SetTime(h, m, s)
	}
	return false
}

func CheckTimeString(target string) bool {
	return NewTime().SetTimeString(target)
}

func (t *Time) SetDelay(h int, m int, s int) bool {
	log.Printf("Time->SetDelay '%d:%02d:%02d'", h, m, s)
	if CheckTime(h, m, s) {
		d, e := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		if e != nil {
			log.Printf("Time->SetDelay ERROR: %s", e.Error())
		} else {
			currentTime := time.Now().In(t.location)
			t.time = currentTime.Add(d)
			t.hours = t.time.Hour()
			t.minutes = t.time.Minute()
			t.seconds = t.time.Second()
			log.Printf("Time->SetDelay set to '%d:%02d:%02d'", t.hours, t.minutes, t.seconds)
			return true
		}
	}
	return false
}

func (t *Time) SetDelayString(delay string) bool {
	log.Printf("Time->SetDelayString '%s'", delay)
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

func CheckDelayString(delay string) bool {
	return NewTime().SetDelayString(delay)
}

func SetCurrentTargetLocation(s string) bool {
	return CurrentTargetTime.SetLocation(s)
}
func SetCurrentTargetTime(h int, m int, s int) bool {
	return CurrentTargetTime.SetTime(h, m, s)
}

func SetCurrentTargetTimeString(s string) bool {
	return CurrentTargetTime.SetTimeString(s)
}

func SetCurrentTargetDelay(h int, m int, s int) bool {
	return CurrentTargetTime.SetDelay(h, m, s)
}

func SetCurrentTargetDelayString(s string) bool {
	return CurrentTargetTime.SetDelayString(s)
}

func SetNextTargetTime() {
	currentTime := time.Now().In(time.UTC)
	SetCurrentTargetLocation("UTC")
	if currentTime.Hour() < 11 {
		SetCurrentTargetTime(11, 15, 0)
	} else if currentTime.Hour() < 14 {
		SetCurrentTargetTime(14, 15, 0)
	} else {
		SetCurrentTargetTime(11, 15, 0)
	}
}
