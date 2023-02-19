package timer

import (
	"log"
)

type TargetListItem struct {
	Object     *Time
	timeString string
	textLabel  string
	alertSound string
	last       struct {
		error error
	}
}

func NewTargetListItem(timeString string, textLabel string, alertSound string) *TargetListItem {
	t := TargetListItem{}
	t.Object = NewTime().SetTimeString(timeString)
	if t.last.error != nil {
		log.Printf("TargetListItem->NewTargetListItem error: %s", t.last.error.Error())
	} else {
		t.timeString = t.Object.TimeString()
		t.textLabel = textLabel
		t.alertSound = alertSound
	}
	return &t
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

func (t *TargetListItem) AlertSound() string {
	return t.alertSound
}
