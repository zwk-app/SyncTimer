package main

import (
	"SyncTimer/app"
	"SyncTimer/app/ui"
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
	appEngine := app.NewAppEngine(ApplicationName, MajorVersion, MinorVersion, BuildNumber, &EmbeddedFS, "res/")
	appEngine.SetTargetJson(appEngine.Config.Target.JsonName).NextTarget()
	if appEngine.Config.Audio.Make {
		appEngine.Audio.GenerateAllAudioFiles(appEngine.AppName())
	} else {
		ui.MainApp(appEngine)
	}
}
