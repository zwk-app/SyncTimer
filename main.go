package main

import (
	"SyncTimer/ttm"
	"SyncTimer/tts"
	"SyncTimer/ui"
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const mainAppName = "SyncTimer"

const embeddedAudioPath = "audio/"

//go:embed audio/*.mp3
var embeddedFS embed.FS

var mainAppVersion string
var mainAppPath string
var verbose bool
var audioPath string
var generateTTS bool
var targetTime string
var targetDelay string

// FixTimezone https://github.com/golang/go/issues/20455
func FixTimezone() {
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}

func main() {
	execFullPath, e := os.Executable()
	if e != nil {
		log.Printf("%s: %s", mainAppName, e.Error())
		os.Exit(1)
	}
	mainAppPath = filepath.Dir(execFullPath)
	mainAppVersion = fmt.Sprintf("%d.%d.%02d", MajorVersion, MinorVersion, BuildNumber)
	tts.SetEmbeddedAudioFS(&embeddedFS, embeddedAudioPath)
	audioPath = mainAppPath + string(os.PathSeparator) + "audio" + string(os.PathSeparator)
	flag.BoolVar(&verbose, "verbose", false, "active stdout logs")
	flag.StringVar(&audioPath, "audioPath", audioPath, "enforce audio local path")
	flag.BoolVar(&generateTTS, "generate-tts", false, "create TTS files")
	flag.StringVar(&targetTime, "time", "", "set target time to <hh[mm[ss]]>")
	flag.StringVar(&targetDelay, "delay", "", "set target time in <[[hh]mm]ss>")
	flag.Parse()
	log.SetPrefix(mainAppName + " v" + mainAppVersion + " ")
	if verbose {
		log.SetOutput(os.Stdout)
	}
	ttm.CurrentTargetTime = ttm.NewTime()
	tts.TextToSpeechEngine = tts.NewTextToSpeech(mainAppName, audioPath, "en")
	if generateTTS {
		tts.GenerateAllAudioFiles(mainAppName, audioPath)
	} else {
		FixTimezone()
		enforceTarget := false
		if targetTime != "" {
			_ = ttm.SetCurrentTargetTimeString(targetTime)
			enforceTarget = true
		}
		if targetDelay != "" {
			_ = ttm.SetCurrentTargetDelayString(targetDelay)
			enforceTarget = true
		}
		ui.MainApp(mainAppName, mainAppVersion, enforceTarget)
	}
}
