package ui

import (
	"SyncTimer/audio"
	"SyncTimer/config"
	"SyncTimer/logs"
	"SyncTimer/timer"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"os"
	"time"
)

var FyneApp fyne.App
var FyneWindow fyne.Window

func TextToSpeechAlert(name string) {
	if config.Alerts().TextToSpeech {
		logs.Debug("MainApp", fmt.Sprintf("TextToSpeechAlert '%s'", name), nil)
		go audio.Play(name)
	}
}

func NotificationAlert(message string) {
	if config.Alerts().Notifications {
		logs.Debug("MainApp", fmt.Sprintf("NotificationAlert '%s'", message), nil)
		go FyneApp.SendNotification(fyne.NewNotification(config.Title(), message))
	}
}

func AlertLoop() {
	logs.Debug("MainApp", "AlertLoop", nil)
	time.Sleep(1500 * time.Millisecond)
	go func() {
		currentCheck := 0
		lastCheck := 0
		lastCheckDiff := 0
		for {
			currentCheck = timer.Engine().RemainingSeconds()
			lastCheckDiff = lastCheck - currentCheck
			if lastCheck < currentCheck {
				logs.Debug("MainApp", fmt.Sprintf("AlertLoop : %08d << %08d (%02d)", lastCheck, currentCheck, lastCheckDiff), nil)
				lastCheck = currentCheck + 1
			}
			if (lastCheckDiff > 1) && (lastCheckDiff < 5) {
				logs.Debug("MainApp", fmt.Sprintf("AlertLoop : %08d <> %08d (%02d)", lastCheck, currentCheck, lastCheckDiff), nil)
				currentCheck = lastCheck - 1
			}
			if currentCheck < lastCheck {
				h, m, s := timer.Engine().RemainingTime()
				if currentCheck >= 0 {
					if currentCheck < 11 {
						if currentCheck == 0 {
							TextToSpeechAlert(config.Alerts().AlarmSound)
						} else if currentCheck%2 == 0 {
							// every 2 sec if T <= 10s
							TextToSpeechAlert(fmt.Sprintf("target-%02d-seconds", s))
						}
					} else if currentCheck < 61 {
						if currentCheck%10 == 0 {
							// every 10 sec if T <= 1m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-seconds", s))
						}
					} else if currentCheck < 310 { // 5m = 300s
						if currentCheck%60 == 0 {
							// every min if T <= 5m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if currentCheck%300 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if currentCheck < 910 { // 15m = 900s
						if currentCheck%300 == 0 {
							// every 5 min if T <= 15m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if currentCheck%900 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if currentCheck < 1810 { // 30m = 1800s
						if currentCheck%600 == 0 {
							// every 10 min if T <= 30m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if currentCheck%1800 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if currentCheck < 10810 { // 3h = 10800s
						if currentCheck%3600 == 0 {
							// every 1 hour if T <= 3h
							TextToSpeechAlert(fmt.Sprintf("target-%02d-hours", h))
						}
					}
				}
				lastCheck = currentCheck
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func MainApp() {
	e := os.Setenv("FYNE_THEME", "dark")
	if e != nil {
		logs.CriticalExit("", "", e)
	}
	FyneApp = app.NewWithID(config.Name())
	config.LoadFyneSettings(FyneApp)
	AlertLoop()
	FyneWindow = FyneApp.NewWindow(config.Title())
	FyneWindow.Resize(fyne.NewSize(320, 540))
	ShowMainWindow()
	FyneApp.Run()
}
