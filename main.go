package main

import (
	"SyncTimer/timer"
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

const embeddedAudioPath = "audio/"

//go:embed audio/*.mp3
var EmbeddedFS embed.FS
var AppName = "SyncTimer"
var AppVersion string
var AppPath string
var TextToSpeechEngine *tts.TextToSpeech
var Timer *timer.TargetTimer
var verbose bool
var audioPath string
var generateTTS bool
var targetTime string
var targetDelay string

// FixTimezone https://github.com/golang/go/issues/20455
func FixTimezone() {
	//goland:noinspection SpellCheckingInspection
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
		log.Printf("%s: %s", AppName, e.Error())
		os.Exit(1)
	}
	AppPath = filepath.Dir(execFullPath)
	AppVersion = fmt.Sprintf("%d.%d.%02d", MajorVersion, MinorVersion, BuildNumber)

	audioPath = AppPath + string(os.PathSeparator) + "audio" + string(os.PathSeparator)
	flag.BoolVar(&verbose, "verbose", false, "active stdout logs")
	flag.StringVar(&audioPath, "audioPath", audioPath, "enforce audio local path")
	flag.BoolVar(&generateTTS, "generate-tts", false, "create TTS files")
	flag.StringVar(&targetTime, "time", "", "set target time to <hh[mm[ss]]>")
	flag.StringVar(&targetDelay, "delay", "", "set target time in <[[hh]mm]ss>")
	flag.Parse()
	log.SetPrefix(AppName + " v" + AppVersion + " ")
	if verbose {
		log.SetOutput(os.Stdout)
	}
	TextToSpeechEngine = tts.NewTextToSpeech(AppName, audioPath, "en")
	TextToSpeechEngine.SetEmbeddedAudioFS(&EmbeddedFS, embeddedAudioPath)
	Timer = timer.NewTargetTimer()

	if generateTTS {
		tts.GenerateAllAudioFiles(AppName, audioPath)
	} else {
		FixTimezone()
		enforceTarget := false
		if targetTime != "" {
			_ = Timer.SetTargetString(targetTime)
			enforceTarget = true
		}
		if targetDelay != "" {
			_ = Timer.SetDelayString(targetDelay)
			enforceTarget = true
		}
		ui.MainApp(AppName, AppVersion, TextToSpeechEngine, Timer, enforceTarget)
	}
}
