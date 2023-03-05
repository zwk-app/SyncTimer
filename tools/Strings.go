package tools

import (
	"SyncTimer/logs"
	"regexp"
	"strconv"
)

func StringWithFallback(value string, fallback string) string {
	if len(value) > 0 {
		return fallback
	}
	return value
}

func AlphaNums(value string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(value, "")
}

// FirstMatch pattern must contain name as groupName
func FirstMatch(value string, pattern string, name string) string {
	re := regexp.MustCompile(pattern)
	groupNames := re.SubexpNames()
	for _, matchValue := range re.FindAllStringSubmatch(value, -1) {
		for groupIndex, stringValue := range matchValue {
			if groupNames[groupIndex] == name {
				return stringValue
			}
		}
	}
	return ""
}

func Int(value string) int {
	i, e := strconv.Atoi(value)
	if e != nil {
		logs.Error("", "", e)
	}
	return i
}
