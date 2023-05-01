package timer

import (
	"SyncTimer/config"
)

type TargetListItem struct {
	Object     *Time
	timeString string
	alarm      struct {
		name  string
		sound string
	}
	last struct {
		error error
	}
}

func NewTargetListItem(timeString string, alarmName string, alarmSound string) *TargetListItem {
	t := new(TargetListItem)
	t.Object = NewTime().SetTimeString(timeString)
	if t.last.error == nil {
		t.timeString = t.Object.TimeString()
		if len(alarmName) > 0 {
			t.alarm.name = alarmName
		} else {
			t.alarm.name = config.DefaultTextLabel
		}
		if len(alarmSound) > 0 {
			t.alarm.sound = alarmSound
		} else {
			t.alarm.sound = config.DefaultAlarmSound
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

func (t *TargetListItem) AlarmName() string {
	return t.alarm.name
}

func (t *TargetListItem) AlarmSound() string {
	return t.alarm.sound
}
