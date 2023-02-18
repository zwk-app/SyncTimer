package tools

import (
	"log"
	"os"
)

func ErrorExit(msg string) {
	log.Printf("ERROR: %s", msg)
	_, _ = os.Stderr.WriteString(msg)
	os.Exit(1)
}
