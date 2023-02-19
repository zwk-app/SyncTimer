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
	FixTimezone()
	appEngine := tools.NewAppEngine(ApplicationName, MajorVersion, MinorVersion, BuildNumber, &EmbeddedFS)
	appEngine.LoadEnvSettings().LoadFileSettings("").LoadArgSettings().SetLogOptions()
	appEngine.Audio.Engine = audio.NewAudioEngine(appEngine.EmbeddedFS, appEngine.Audio.EmbeddedPath, appEngine.Audio.LocalPath, "en")
	appEngine.Timer.List = timer.NewTargetList().LoadJsonURL(appEngine.Timer.ListURL).LoadJsonFile(appEngine.Timer.ListFile)
	appEngine.Timer.Engine = timer.NewTargetTimer()
	appEngine.NextTarget()
	if appEngine.Audio.GenerateTTS {
		appEngine.Audio.Engine.GenerateAllAudioFiles(appEngine.AppName())
	} else {
		ui.MainApp(appEngine)
	}
}
