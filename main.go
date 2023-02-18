package main

import (
	"SyncTimer/timer"
	"SyncTimer/tools"
	"SyncTimer/tts"
	"SyncTimer/ui"
	"embed"
	"os/exec"
	"strings"
	"time"
)

//go:embed audio/*.mp3
var EmbeddedFS embed.FS

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
	appEngine := tools.NewAppEngine(ApplicationName, MajorVersion, MinorVersion, BuildNumber, &EmbeddedFS)
	appEngine.LoadEnvSettings().LoadFileSettings("").LoadArgSettings().SetLogOptions()
	appEngine.TextToSpeech.Object = tts.NewTextToSpeech(appEngine.Name(), appEngine.TextToSpeech.LocalPath, "en")
	appEngine.TextToSpeech.Object.SetEmbeddedAudioFS(&EmbeddedFS, appEngine.TextToSpeech.EmbeddedPath)
	appEngine.Timer.Object = timer.NewTargetTimer()

	if appEngine.TextToSpeech.GenerateTTS {
		tts.GenerateAllAudioFiles(appEngine.Name(), appEngine.TextToSpeech.LocalPath)
	} else {
		FixTimezone()
		if appEngine.Timer.TargetTime != "" {
			_ = appEngine.Timer.Object.SetTargetString(appEngine.Timer.TargetTime)
			appEngine.Timer.EnforceTarget = true
		}
		if appEngine.Timer.TargetDelay != "" {
			_ = appEngine.Timer.Object.SetDelayString(appEngine.Timer.TargetDelay)
			appEngine.Timer.EnforceTarget = true
		}
		ui.MainApp(appEngine)
	}
}
