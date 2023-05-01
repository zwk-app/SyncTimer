package timer

import (
	"fmt"
	"strconv"
)

func CheckDelay(h int, m int, s int) bool {
	return CheckTime(h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func CheckDelayString(delayString string) bool {
	h, m, s, e := TimeFromString(delayString)
	if e != nil {
		return false
	}
	return CheckDelay(h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func StringFromDelay(h int, m int, s int) string {
	return fmt.Sprintf("%02d%02d%02d", h, m, s)
}

func DelayFromString(delayString string) (int, int, int, error) {
	targetLen := len(delayString)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch targetLen {
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
	if !hasError {
		if CheckDelay(h, m, s) {
			return h, m, s, nil
		}
	}
	return 0, 0, 0, fmt.Errorf("invalid delay string '%s'", delayString)
}
