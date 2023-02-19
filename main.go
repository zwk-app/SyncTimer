package main

import (
	"SyncTimer/audio"
	"SyncTimer/timer"
	"SyncTimer/tools"
	"SyncTimer/ui"
	"embed"
	"os/exec"
	"strings"
	"time"
)

//go:embed res/audio/*.mp3
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
	appEngine.Audio.Object = audio.NewAudioEngine(appEngine.EmbeddedFS, appEngine.Audio.EmbeddedPath, appEngine.Audio.LocalPath, "en")
	appEngine.Timer.Object = timer.NewTargetTimer()

	if appEngine.Audio.GenerateTTS {
		appEngine.Audio.Object.GenerateAllAudioFiles(appEngine.Name())
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
