package timer

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"unicode"
)

type Time struct {
	location     *time.Location
	locationName string
	time         time.Time
	hours        int
	minutes      int
	seconds      int
	last         struct {
		error error
	}
}

//goland:noinspection ALL
const LocalLocationName = "Local Time"

//goland:noinspection ALL
const UtcLocationName = "UTC"

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

func (t *Time) Error() error {
	lastError := t.last.error
	t.last.error = nil
	return lastError
}
func (t *Time) Location() *time.Location {
	return t.location
}

func (t *Time) LocationName() string {
	return t.locationName
}

func (t *Time) Time() time.Time {
	return t.time
}

func (t *Time) TimeString() string {
	return StringFromTime(t.hours, t.minutes, t.seconds)
}

func (t *Time) TextString() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.hours, t.minutes, t.seconds)
}

func (t *Time) SetLocation(location *time.Location) *Time {
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
	}
	return t
}

func (t *Time) SetLocationName(location string) *Time {
	if t.locationName != location {
		switch location {
		case "UTC":
			t.SetLocation(time.UTC)
		case LocalLocationName:
			t.SetLocation(time.Local)
		default:
			t.SetLocation(time.Local)
		}
	}
	return t
}

func (t *Time) Set(new time.Time) *Time {
	t.time = new.In(t.location)
	t.hours = t.time.Hour()
	t.minutes = t.time.Minute()
	t.seconds = t.time.Second()
	return t
}

func (t *Time) SetTime(h int, m int, s int) *Time {
	if !CheckTime(h, m, s) {
		t.last.error = fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
	} else {
		currentTime := time.Now().In(t.location)
		tmpTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), h, m, s, 0, t.location)
		if tmpTime.Before(currentTime) {
			tmpTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, h, m, s, 0, t.location)
		}
		t.Set(tmpTime)
	}
	return t
}

func (t *Time) SetTimeString(timeString string) *Time {
	currentLocation := t.location
	if len(timeString) > 0 {
		if unicode.IsLetter(rune(timeString[0])) {
			switch timeString[0] {
			case 'U':
				if currentLocation != time.UTC {
					t.SetLocation(time.UTC)
				}
			case 'L':
				if currentLocation != time.Local {
					t.SetLocation(time.Local)
				}
			default:
				t.last.error = fmt.Errorf("invalid time string '%s'", timeString)
				log.Printf("Time->SetTimeString error: %s", t.last.error.Error())
				return t
			}
			timeString = timeString[1:]
		}
		h, m, s, e := TimeFromString(timeString)
		if e != nil {
			t.last.error = e
			log.Printf("Time->SetTimeString error: %s", t.last.error.Error())
			return t
		}
		t.SetTime(h, m, s)
		if t.last.error != nil {
			log.Printf("Time->SetTimeString error: %s", t.last.error.Error())
			return t
		}
		if currentLocation != t.location {
			_ = t.SetLocation(currentLocation)
		}
		return t
	}
	t.last.error = fmt.Errorf("invalid time string '%s'", timeString)
	log.Printf("Time->SetTimeString error: %s", t.last.error.Error())
	return t
}

func CheckTime(h int, m int, s int) bool {
	if h >= 0 && h <= 23 && m >= 0 && m <= 59 && s >= 0 && s <= 59 {
		return true
	}
	return false
}

func CheckTimeString(timeString string) bool {
	h, m, s, e := TimeFromString(timeString)
	if e != nil {
		return false
	}
	return CheckTime(h, m, s)
}

func StringFromTime(h int, m int, s int) string {
	return fmt.Sprintf("%02d%02d%02d", h, m, s)
}

func TimeFromString(timeString string) (int, int, int, error) {
	targetLen := len(timeString)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch targetLen {
	case 1, 2:
		h, he = strconv.Atoi(timeString)
		if he != nil {
			hasError = true
		}
	case 3:
		h, he = strconv.Atoi(timeString[0:1])
		m, me = strconv.Atoi(timeString[1:3])
		if he != nil || me != nil {
			hasError = true
		}
	case 4:
		h, he = strconv.Atoi(timeString[0:2])
		m, me = strconv.Atoi(timeString[2:4])
		if he != nil || me != nil {
			hasError = true
		}
	case 5:
		h, he = strconv.Atoi(timeString[0:1])
		m, me = strconv.Atoi(timeString[1:3])
		s, se = strconv.Atoi(timeString[3:5])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	case 6:
		h, he = strconv.Atoi(timeString[0:2])
		m, me = strconv.Atoi(timeString[2:4])
		s, se = strconv.Atoi(timeString[4:6])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	}
	if !hasError {
		if CheckTime(h, m, s) {
			return h, m, s, nil
		}
	}
	return 0, 0, 0, fmt.Errorf("invalid time string '%s'", timeString)
}
