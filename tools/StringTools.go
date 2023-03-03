package tools

import "regexp"

func AlphaNums(value string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(value, "")
}
