package main

import (
	"SyncTimer/audio"
	"SyncTimer/config"
	"SyncTimer/resources"
	"SyncTimer/ui"
	"embed"
	"fmt"
	"github.com/zwk-app/zwk-tools/logs"
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
	if config.Logs().Verbose {
		logs.SetLevelDebug()
	} else {
		logs.SetLevelInfo()
	}
	if len(config.Logs().FileName) > 0 {
		logs.SetFileName(config.Logs().FileName)
	}
	logs.Info("Main", fmt.Sprintf("CurrentConfig\n%s", config.ToString()), nil)
	if config.Config().Audio.Make {
		audio.GenerateAll(config.Name())
	} else {
		ui.MainApp()
	}
}
