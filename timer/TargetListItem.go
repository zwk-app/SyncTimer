package timer

import (
	"SyncTimer/config"
)

type TargetListItem struct {
	Object     *Time
	timeString string
	textLabel  string
	alarmSound string
	last       struct {
		error error
	}
}

func NewTargetListItem(timeString string, textLabel string, alertSound string) *TargetListItem {
	t := new(TargetListItem)
	t.Object = NewTime().SetTimeString(timeString)
	if t.last.error == nil {
		t.timeString = t.Object.TimeString()
		if len(textLabel) > 0 {
			t.textLabel = textLabel
		} else {
			t.textLabel = config.DefaultTextLabel
		}
		if len(alertSound) > 0 {
			t.alarmSound = alertSound
		} else {
			t.alarmSound = config.DefaultAlarmSound
		}
	}
	return t
}

func (t *TargetListItem) Error() error {
	lastError := t.last.error
	t.last.error = nil
	return lastError
}

func (t *TargetListItem) Time() (int, int, int) {
	h, m, s, _ := TimeFromString(t.timeString)
	return h, m, s
}

func (t *TargetListItem) TimeString() string {
	return t.timeString
}

func (t *TargetListItem) TextLabel() string {
	return t.textLabel
}

func (t *TargetListItem) AlarmSound() string {
	return t.alarmSound
}
