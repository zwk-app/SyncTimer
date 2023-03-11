package timer

import (
	"SyncTimer/config"
	"fmt"
	"github.com/zwk-app/go-tools/logs"
	"github.com/zwk-app/go-tools/tools"
	"strconv"
	"time"
)

type TargetEngine struct {
	current    *Time
	target     *Time
	List       *TargetList
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

var engine *TargetEngine

func Engine() *TargetEngine {
	if engine == nil {
		logs.Debug("TargetEngine", fmt.Sprintf("Init"), nil)
		engine = NewTargetEngine()
	}
	return engine
}

func NewTargetEngine() *TargetEngine {
	t := TargetEngine{}
	t.current = NewTime()
	t.target = NewTime()
	t.textLabel = config.DefaultTextLabel
	t.alertSound = config.DefaultAlarmSound
	t.remaining.timeInSeconds = 0
	t.remaining.hours = 0
	t.remaining.minutes = 0
	t.remaining.seconds = 0
	ch := make(chan time.Time)
	t.remainingTimeLoop(ch)
	t.currentTimeLoop(ch)
	return &t
}

func SetError(e error) {
	Engine().last.error = e
}

func HasError() bool {
	return Engine().last.error != nil
}

func GetError() error {
	lastError := Engine().last.error
	SetError(nil)
	return lastError
}

func (t *TargetEngine) currentTimeLoop(ch chan time.Time) {
	logs.Debug("TargetEngine", fmt.Sprintf("CurrentTimeLoop"), nil)
	go func() {
		for {
			current := time.Now()
			t.current.Set(current)
			ch <- current
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func (t *TargetEngine) remainingTimeLoop(ch chan time.Time) {
	logs.Debug("TargetEngine", fmt.Sprintf("RemainginTimeLoop"), nil)
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

func (t *TargetEngine) TextLabel() string {
	return t.textLabel
}

func (t *TargetEngine) AlertSound() string {
	return t.alertSound
}

func (t *TargetEngine) LocationName() string {
	return t.target.LocationName()
}

func (t *TargetEngine) CurrentTimeString() string {
	return t.current.TimeString()
}

func (t *TargetEngine) CurrentTextString() string {
	return t.current.TextString()
}

func (t *TargetEngine) TargetTimeString() string {
	return t.target.TimeString()
}

func (t *TargetEngine) TargetTextString() string {
	return t.target.TextString()
}

func (t *TargetEngine) RemainingString() string {
	if t.remaining.hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", t.remaining.hours, t.remaining.minutes, t.remaining.seconds)
	}
	if t.remaining.minutes > 0 {
		return fmt.Sprintf("%d:%02d", t.remaining.minutes, t.remaining.seconds)
	}
	return fmt.Sprintf("%d", t.remaining.seconds)
}

func (t *TargetEngine) RemainingTime() (int, int, int) {
	return t.remaining.hours, t.remaining.minutes, t.remaining.seconds
}

func (t *TargetEngine) RemainingSeconds() int {
	return t.remaining.timeInSeconds
}

func (t *TargetEngine) SetTextLabel(textLabel string) *TargetEngine {
	if len(textLabel) > 0 {
		t.textLabel = textLabel
	} else {
		t.textLabel = config.DefaultTextLabel
	}
	return t
}

func (t *TargetEngine) SetAlarmSound(alarmSound string) *TargetEngine {
	if len(alarmSound) > 0 {
		t.alertSound = alarmSound
	} else {
		t.alertSound = config.DefaultAlarmSound
	}
	return t
}

func (t *TargetEngine) SetLocationName(locationName string) *TargetEngine {
	targetLocation := t.target.LocationName()
	currentLocation := t.current.LocationName()
	if targetLocation != locationName || currentLocation != locationName {
		t.target.SetLocationName(locationName)
		t.current.SetLocationName(locationName)
	}
	return t
}

func (t *TargetEngine) SetDelay(h int, m int, s int) *TargetEngine {
	logs.Debug("TargetEngine", fmt.Sprintf("SetDelay: %02d!%02d:%02d", h, m, s), nil)
	if !CheckTime(h, m, s) {
		t.last.error = fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
		logs.Debug("", "", t.last.error)
	} else {
		d, e := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		if e != nil {
			t.last.error = e
			logs.Debug("", "", t.last.error)
		} else {
			currentTime := time.Now().In(t.target.location)
			t.target.Set(currentTime.Add(d))
			if t.last.error != nil {
				logs.Debug("", "", t.last.error)
			}
		}
	}
	return t
}

func (t *TargetEngine) SetDelayString(delayString string) *TargetEngine {
	logs.Debug("TargetEngine", fmt.Sprintf("SetDelayString: %s", delayString), nil)
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
		logs.Debug("", "", t.last.error)
	} else {
		t.SetDelay(h, m, s)
		if t.last.error != nil {
			logs.Debug("", "", t.last.error)
		}
	}
	return t
}

func (t *TargetEngine) SetTarget(value *Time) *TargetEngine {
	logs.Debug("TargetEngine", fmt.Sprintf("SetTarget: %T", value), nil)
	t.target = value
	return t
}

func (t *TargetEngine) SetTargetTime(h int, m int, s int) *TargetEngine {
	logs.Debug("TargetEngine", fmt.Sprintf("SetTargetTime: %02d!%02d:%02d", h, m, s), nil)
	t.target.SetTime(h, m, s)
	if t.last.error != nil {
		logs.Debug("", "", t.last.error)
	}
	return t
}
func (t *TargetEngine) SetTargetString(value string) *TargetEngine {
	logs.Debug("TargetEngine", fmt.Sprintf("SetTargetString: %s", value), nil)
	t.target.SetTimeString(value)
	if t.last.error != nil {
		logs.Debug("", "", t.last.error)
	}
	return t
}

func SetAlarmSound(value string) {
	Engine().SetAlarmSound(tools.Fallback(value, config.DefaultAlarmSound))
}

func SetTextLabel(value string) {
	Engine().SetTextLabel(tools.Fallback(value, config.DefaultTextLabel))
}

func SetTargetTime(timeString string, label string, alarm string) {
	logs.Debug("AudioEngine", fmt.Sprintf("SetTargetTime: '%s'", timeString), nil)
	Engine().SetTargetString(timeString)
	if HasError() {
		logs.Error("AudioEngine", fmt.Sprintf("SetTargetTime: '%s'", timeString), nil)
	} else {
		SetTextLabel(label)
		SetAlarmSound(alarm)
	}
}

func SetTargetDelay(delayString string, label string, alarm string) {
	logs.Debug("AudioEngine", fmt.Sprintf("SetTargetDelay: '%s'", delayString), nil)
	Engine().SetDelayString(delayString)
	if HasError() {
		logs.Error("AudioEngine", fmt.Sprintf("SetTargetDelay: '%s'", delayString), nil)
	} else {
		SetTextLabel(label)
		SetAlarmSound(alarm)
	}
}

func SetTargetJson(value string) {
	logs.Debug("AudioEngine", fmt.Sprintf("SetTargetJson: '%s'", value), nil)
	r := Engine()
	if r.List == nil {
		r.List = NewTargetList()
	}
	r.List.LoadJson(value)
}

func NextTarget() {
	logs.Debug("AudioEngine", fmt.Sprintf("NextTarget"), nil)
	if config.Config().Target.Delay != "" {
		SetTargetDelay(config.Config().Target.Delay, "", "")
	} else if config.Config().Target.Time != "" {
		SetTargetTime(config.Config().Target.Time, "", "")
	} else {
		next := Engine().List.NextTargetListItem()
		SetTargetTime(next.timeString, next.textLabel, next.alarmSound)
	}
}
