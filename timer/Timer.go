package timer

import (
	"SyncTimer/config"
	"fmt"
	"github.com/zwk-app/go-tools/logs"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type TimerReceiver struct {
	current   *Time
	target    *Time
	List      *TargetList
	remaining struct {
		timeInSeconds int64
		hours         int
		minutes       int
		seconds       int
	}
	alarm struct {
		name  string
		sound string
	}
	last struct {
		error error
	}
}

var timer *TimerReceiver

func Timer() *TimerReceiver {
	if timer == nil {
		logs.Debug("Timer", fmt.Sprintf("Init"), nil)
		timer = new(TimerReceiver)
		timer.current = NewTime()
		timer.target = NewTime()
		timer.alarm.name = config.DefaultTextLabel
		timer.alarm.sound = config.DefaultAlarmSound
		timer.remaining.timeInSeconds = 0
		timer.remaining.hours = 0
		timer.remaining.minutes = 0
		timer.remaining.seconds = 0
		ch := make(chan time.Time)
		timer.remainingTimeLoop(ch)
		timer.currentTimeLoop(ch)
	}
	return timer
}

func SetError(e error) {
	Timer().last.error = e
}

func HasError() bool {
	return Timer().last.error != nil
}

//goland:noinspection GoUnusedExportedFunction
func GetError() error {
	if HasError() {
		lastError := Timer().last.error
		SetError(nil)
		return lastError
	}
	return nil
}

func (t *TimerReceiver) currentTimeLoop(ch chan time.Time) {
	logs.Debug("Timer", fmt.Sprintf("CurrentTimeLoop"), nil)
	go func() {
		for t != nil {
			current := time.Now()
			t.current.Set(current)
			ch <- current
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func (t *TimerReceiver) remainingTimeLoop(ch chan time.Time) {
	logs.Debug("Timer", fmt.Sprintf("RemainginTimeLoop"), nil)
	go func() {
		for current := range ch {
			d := t.target.time.Sub(current).Round(time.Second)
			t.remaining.timeInSeconds = int64(d / time.Second)
			h := d / time.Hour
			d -= h * time.Hour
			m := d / time.Minute
			d -= m * time.Minute
			s := d / time.Second
			t.remaining.hours = int(h)
			t.remaining.minutes = int(m)
			t.remaining.seconds = int(s)
		}
	}()
}

//goland:noinspection GoUnusedExportedFunction
func AlarmName() string {
	return Timer().alarm.name
}

//goland:noinspection GoUnusedExportedFunction
func AlarmSound() string {
	return Timer().alarm.sound
}

//goland:noinspection GoUnusedExportedFunction
func LocationName() string {
	return Timer().target.LocationName()
}

//goland:noinspection GoUnusedExportedFunction
func CurrentTimeString() string {
	return Timer().current.TimeString()
}

//goland:noinspection GoUnusedExportedFunction
func CurrentTimeText() string {
	return Timer().current.TextString()
}

//goland:noinspection GoUnusedExportedFunction
func TargetTimeString() string {
	return Timer().target.TimeString()
}

//goland:noinspection GoUnusedExportedFunction
func TargetTimeText() string {
	return Timer().target.TextString()
}

//goland:noinspection GoUnusedExportedFunction
func RemainingSeconds() int64 {
	return Timer().remaining.timeInSeconds
}

//goland:noinspection GoUnusedExportedFunction
func RemainingTime() (int, int, int) {
	t := Timer()
	return t.remaining.hours, t.remaining.minutes, t.remaining.seconds
}

//goland:noinspection GoUnusedExportedFunction
func RemainingTimeText() string {
	t := Timer()
	if t.remaining.hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", t.remaining.hours, t.remaining.minutes, t.remaining.seconds)
	}
	if t.remaining.minutes > 0 {
		return fmt.Sprintf("%d:%02d", t.remaining.minutes, t.remaining.seconds)
	}
	return fmt.Sprintf("%d", t.remaining.seconds)
}

//goland:noinspection GoUnusedExportedFunction
func SetAlarmName(alarmName string) {
	if len(alarmName) > 0 {
		Timer().alarm.name = alarmName
	} else {
		Timer().alarm.name = config.DefaultTextLabel
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetAlarmSound(alarmSound string) {
	if len(alarmSound) > 0 {
		Timer().alarm.sound = alarmSound
	} else {
		Timer().alarm.sound = config.DefaultAlarmSound
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetLocationName(locationName string) {
	t := Timer()
	targetLocation := t.target.LocationName()
	currentLocation := t.current.LocationName()
	if targetLocation != locationName || currentLocation != locationName {
		t.target.SetLocationName(locationName)
		t.current.SetLocationName(locationName)
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetTargetDelay(h int, m int, s int) {
	logs.Debug("Timer", fmt.Sprintf("SetTargetDelay: %02d!%02d:%02d", h, m, s), nil)
	t := Timer()
	if !CheckDelay(h, m, s) {
		t.last.error = fmt.Errorf("invalid delay '%02d:%02d:%02d'", h, m, s)
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
}

//goland:noinspection GoUnusedExportedFunction
func SetTargetDelayString(delayString string) {
	logs.Debug("Timer", fmt.Sprintf("SetTargetDelayString: %s", delayString), nil)
	h, m, s, e := DelayFromString(delayString)
	if e == nil {
		SetTargetDelay(h, m, s)
	} else {
		SetError(e)
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetTarget(value *Time) {
	logs.Debug("Timer", fmt.Sprintf("SetTarget: %T", value), nil)
	Timer().target = value
}

//goland:noinspection GoUnusedExportedFunction
func SetTargetTime(h int, m int, s int) {
	logs.Debug("Timer", fmt.Sprintf("SetTargetTime: %02d!%02d:%02d", h, m, s), nil)
	Timer().target.SetTime(h, m, s)
	if HasError() {
		logs.Debug("", "", Timer().last.error)
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetTargetTimeString(value string) {
	logs.Debug("Timer", fmt.Sprintf("SetTargetTimeString: '%s'", value), nil)
	Timer().target.SetTimeString(value)
	if HasError() {
		logs.Debug("", "", Timer().last.error)
	}
}

func SetTargetJson(value string) {
	logs.Debug("Timer", fmt.Sprintf("SetTargetJson: '%s'", value), nil)
	r := Timer()
	if r.List == nil {
		r.List = NewTargetList()
	}
	r.List.LoadJson(value)
}

func NextTarget() {
	logs.Debug("Timer", fmt.Sprintf("NextTarget"), nil)
	if config.Config().Target.Delay != "" {
		SetTargetDelayString(config.Config().Target.Delay)
	} else if config.Config().Target.Time != "" {
		SetTargetTimeString(config.Config().Target.Time)
	} else {
		next := Timer().List.NextTargetListItem()
		SetTargetTimeString(next.timeString)
		SetAlarmName(next.alarm.name)
		SetAlarmSound(next.alarm.sound)
	}
}
