package main

import (
	"SyncTimer/audio"
	"SyncTimer/config"
	"SyncTimer/logs"
	"SyncTimer/resources"
	"SyncTimer/timer"
	"SyncTimer/ui"
	"embed"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

//go:embed res/icon.png res/images/*.svg res/audio/*.mp3
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
	resources.SetEmbedded(&EmbeddedFS)
	config.SetAppInfo(ApplicationName, MajorVersion, MinorVersion, BuildNumber)
	config.LoadEnvironment()
	_ = config.LoadFile("")
	config.LoadArguments()
	if currentConfig, e := config.ToJson(); e == nil {
		logs.Debug("Main", fmt.Sprintf("CurrentConfig: %s", currentConfig), nil)
	}
	timer.SetTargetJson(config.Target().JsonName)
	timer.NextTarget()
	if config.Config().Audio.Make {
		audio.GenerateAll(config.Name())
	} else {
		ui.MainApp()
	}

}
