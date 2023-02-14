package timer

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

const LocalLocationName = "Local Time"

func NewTime() *Time {
	return &Time{
		location:     time.Local,
		locationName: LocalLocationName,
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

func (t *Time) GetTimeString() string {
	return fmt.Sprintf("%d:%02d:%02d", t.hours, t.minutes, t.seconds)
}

func (t *Time) SetLocation(location *time.Location) bool {
	log.Printf("Time->SetLocation '%s'", location)
	if t.location != location {
		switch location {
		case time.UTC:
			t.location = time.UTC
			t.locationName = "UTC"
		case time.Local:
			t.location = time.Local
			t.locationName = LocalLocationName
		default:
			t.location = time.Local
			t.locationName = LocalLocationName
		}
		t.Set(t.time.In(location))
		return true
	}
	return false
}

func (t *Time) SetLocationName(location string) bool {
	log.Printf("Time->SetLocationName '%s'", location)
	if t.locationName != location {
		switch location {
		case "UTC":
			t.SetLocation(time.UTC)
		case LocalLocationName:
			t.SetLocation(time.Local)
		default:
			t.SetLocation(time.Local)
		}
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

func (t *Time) Set(new time.Time) {
	t.time = new.In(t.location)
	t.hours = t.time.Hour()
	t.minutes = t.time.Minute()
	t.seconds = t.time.Second()
}

func (t *Time) SetTime(h int, m int, s int) bool {
	log.Printf("Time->SetTime '%d:%02d:%02d'", h, m, s)
	if CheckTime(h, m, s) {
		currentTime := time.Now().In(t.location)
		tmpTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), h, m, s, 0, t.location)
		if tmpTime.Before(currentTime) {
			tmpTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, h, m, s, 0, t.location)
		}
		t.Set(tmpTime)
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
